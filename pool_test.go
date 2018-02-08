package pool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var (
	numWorkers int64 = 4
	t10ms            = time.Duration(10) * time.Millisecond
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Test page")
}

func testHandlerWithTimeout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(t10ms)
	fmt.Fprint(w, "Test page with timeout")
}

func TestQueue(t *testing.T) {
	queue := new(taskQueue)
	testTask := &Task{
		Hostname: "",
	}
	queue.put(testTask)
	if queue.len != 1 {
		t.Errorf("Got %v error, want %v", queue.len, 1)
	}
	queue.put(testTask)
	_, ok := queue.get()
	if !ok {
		t.Errorf("Got %v error, want %v", ok, true)
	}
}

func TestNoServer(t *testing.T) {
	p := New(numWorkers)
	if p.numWorkers != numWorkers {
		t.Errorf("Got %v numWorkers, want %v", p.numWorkers, numWorkers)
	}
	err := p.Add("", nil)
	if err != errEmptyTarget {
		t.Errorf("Got %v error, want %v", err, errEmptyTarget)
	}
	err = p.Add(":", nil)
	if err == nil {
		t.Errorf("Got %v error, want hostname error", err)
	}
	if !p.poolIsRunning() {
		t.Errorf("Got %v error, want true", p.poolIsRunning())
	}
	err = p.Add("http://127.0.0.1:80/", nil)
	if err != nil {
		t.Errorf("Got %v error, want nil error", err)
	}
	task := <-p.ResultChan
	if task.Error == nil {
		t.Error("Got nil error, want net error")
	}
	if p.GetAddedTasks() != 1 {
		t.Errorf("Wrong input jobs. Want 1, got %v", p.GetAddedTasks())
	}
	if p.getFreeWorkers() != numWorkers {
		t.Errorf("Wrong free workers. Want %v, got %v", numWorkers, p.getFreeWorkers())
	}
	proxy, _ := url.Parse("http://127.0.0.1:80/")
	err = p.Add("http://127.0.0.1:80/", proxy)
	if err != nil {
		t.Errorf("Got %v error, want nil error", err)
	}
	task = <-p.ResultChan
	if task.Error == nil {
		t.Error("Got nil error, want net error")
	}
	p.Quit()
	if p.poolIsRunning() {
		t.Errorf("Got %v error, want false", p.poolIsRunning())
	}
}

func TestWithServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	p := New(numWorkers)
	_ = p.Add(ts.URL, nil)
	task, ok := <-p.ResultChan
	if !ok {
		t.Error("Channel is closed'")
	}
	if string(task.Body) != "Test page" {
		t.Errorf("Got %v error, want 'Test page'", string(task.Body))
	}
	p.EndWaitingTasks()
	err := p.Add(ts.URL, nil)
	if err != errNotWait {
		t.Errorf("Got %v error, want %v", err, errNotWait)
	}
	if !p.isCompleteJobs() {
		t.Errorf("Got %v error, want %v", false, true)
	}
}

func TestWithTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandlerWithTimeout))
	defer ts.Close()

	p := New(1)
	p.SetTimeout(100)
	err := p.Add(ts.URL, nil)
	if err != nil {
		t.Errorf("Got %v error, want %v", err, nil)
	}
	task, ok := <-p.ResultChan
	if !ok {
		t.Error("Channel is closed'")
	}
	if string(task.Body) != "Test page with timeout" {
		t.Errorf("Got %v error, want 'Test page with timeout'", string(task.Body))
	}
	p.timeout = time.Duration(1) * time.Millisecond
	err = p.Add(ts.URL, nil)
	if err != nil {
		t.Errorf("Got %v error, want %v", err, nil)
	}
	task, ok = <-p.ResultChan
	if !ok {
		t.Error("Channel is closed'")
	}
	if task.Error == nil {
		t.Errorf("Got no error, want %v", task.Error)
	}
	_ = p.Add(ts.URL, nil)
	p.Quit()
	if p.GetCompletedTasks() != 2 {
		t.Errorf("Got %v error, want %v", p.GetCompletedTasks(), 2)
	}
	task, ok = <-p.ResultChan
	if ok {
		t.Errorf("Got %v error, want %v", ok, false)
	}
	if task != nil {
		t.Errorf("Got %v error, want %v", task, nil)
	}
	if p.poolIsRunning() {
		t.Errorf("Got %v error, want %v", !p.poolIsRunning(), false)
	}
	err = p.Add(ts.URL, nil)
	if err != errNotRun {
		t.Errorf("Got %v error, want %v", err, errNotRun)
	}
}

func TestQuitTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandlerWithTimeout))
	defer ts.Close()

	p := New(1)
	p.SetTimeout(30)
	p.SetQuitTimeout(2)
	_ = p.Add(ts.URL, nil)
	_ = p.Add(ts.URL, nil)
	for range p.ResultChan {
	}
	if p.GetCompletedTasks() != 2 {
		t.Errorf("Got %v error, want %v", p.GetCompletedTasks(), 1)
	}
}

func TestWaitingTasks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	p := New(1)
	_ = p.Add(ts.URL, nil)
	_ = p.Add(ts.URL, nil)
	p.EndWaitingTasks()
	for range p.ResultChan {
	}
	if p.GetCompletedTasks() != 2 {
		t.Errorf("Got %v error, want %v", p.GetCompletedTasks(), 1)
	}
}

func BenchmarkAccumulate(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	p := New(numWorkers)
	n := b.N
	for i := 0; i < n; i++ {
		err := p.Add(ts.URL, nil)
		if err != nil {
			b.Errorf("Got %v error, want %v", err, nil)
		}
	}
	for i := 0; i < n; i++ {
		task := <-p.ResultChan
		err := task.Error
		if err != nil {
			b.Error("Error", err)
		}
	}
}

func BenchmarkParallel(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	p := New(numWorkers)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := p.Add(ts.URL, nil)
			if err != nil {
				b.Errorf("Got %v error, want %v", err, nil)
			}
		}
	})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := <-p.ResultChan
			err := task.Error
			if err != nil {
				b.Error("Error", err)
			}
		}
	})
}
