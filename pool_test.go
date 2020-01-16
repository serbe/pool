package pool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	numWorkers int64 = 4
	t10ms            = time.Duration(10) * time.Millisecond
	t30ms            = time.Duration(30) * time.Millisecond
)

func testHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Test page")
}

func testHandlerWithTimeout(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(t10ms)
	fmt.Fprint(w, "Test page with timeout")
}

// func TestQueue(t *testing.T) {
// 	queue := newRingQueue()
// 	testTask := Task{
// 		Hostname: "",
// 	}
// 	queue.put(testTask)
// 	if queue.cnt != 1 {
// 		t.Errorf("Got %v error, want %v", queue.cnt, 1)
// 	}
// 	queue.put(testTask)
// 	_, ok := queue.get()
// 	if !ok {
// 		t.Errorf("Got %v error, want %v", ok, true)
// 	}
// }

func TestClosedInputTaskChanByTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	p := New(2)
	p.SetQuitTimeout(10)
	err := p.Add(ts.URL, "")
	if err != nil {
		t.Errorf("Got %v error, want %v", nil, errNotRun)
	}
	time.Sleep(t30ms)
	err = p.Add(ts.URL, "")
	if err == nil {
		t.Errorf("Got %v error, want %v", nil, errNotRun)
	}
}

func TestNoServer(t *testing.T) {
	p := New(numWorkers)
	if p.numWorkers != numWorkers {
		t.Errorf("Got %v numWorkers, want %v", p.numWorkers, numWorkers)
	}
	err := p.Add("", "")
	if err != errEmptyTarget {
		t.Errorf("Got %v error, want %v", err, errEmptyTarget)
	}
	if !p.poolIsRunning() {
		t.Errorf("Got %v error, want true", p.poolIsRunning())
	}
	err = p.Add(":", "")
	if err != nil {
		t.Errorf("Got %v error, want nil error", err)
	}
	task := <-p.ResultChan
	if task.Error == nil {
		t.Error("Got nil error, want net error")
	}
	err = p.Add("http://127.0.0.1:80/", ":")
	if err != nil {
		t.Errorf("Got %v error, want nil error", err)
	}
	task = <-p.ResultChan
	if task.Error == nil {
		t.Error("Got nil error, want net error")
	}
	err = p.Add("http://127.0.0.1:80/", "")
	if err != nil {
		t.Errorf("Got %v error, want nil error", err)
	}
	task = <-p.ResultChan
	if task.Error == nil {
		t.Error("Got nil error, want net error")
	}
	if p.GetAddedTasks() != 3 {
		t.Errorf("Wrong input jobs. Want 1, got %v", p.GetAddedTasks())
	}
	if p.getFreeWorkers() != numWorkers {
		t.Errorf("Wrong free workers. Want %v, got %v", numWorkers, p.getFreeWorkers())
	}
	err = p.Add("http://127.0.0.1:80/", "http://127.0.0.1:80/")
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
	_ = p.Add(ts.URL, "")
	task, ok := <-p.ResultChan
	if !ok {
		t.Error("Channel is closed'")
	}
	if string(task.Body) != "Test page" {
		t.Errorf("Got %v error, want 'Test page'", string(task.Body))
	}
	p.EndWaitingTasks()
	err := p.Add(ts.URL, "")
	if err != errNotWait {
		t.Errorf("Got %v error, want %v", err, errNotWait)
	}
	if p.GetCompletedTasks() == 0 {
		t.Errorf("Got %v error, want %v", 0, 1)
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
	err := p.Add(ts.URL, "")
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
	err = p.Add(ts.URL, "")
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
	_ = p.Add(ts.URL, "")
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
	err = p.Add(ts.URL, "")
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
	_ = p.Add(ts.URL, "")
	_ = p.Add(ts.URL, "")
	if p.GetCompletedTasks() != 0 {
		t.Errorf("Got %v error, want %v", p.GetCompletedTasks(), 0)
	}
}

func TestWaitingTasks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	p := New(1)
	_ = p.Add(ts.URL, "")
	_ = p.Add(ts.URL, "")
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
	b.ResetTimer()

	p := New(numWorkers)
	n := b.N
	for i := 0; i < n; i++ {
		err := p.Add(ts.URL, "")
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
	b.ResetTimer()

	p := New(numWorkers)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := p.Add(ts.URL, "")
			if err != nil {
				b.Errorf("Got %v error, want %v", err, nil)
			}
		}
	})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := <-p.ResultChan
			_ = task.Error
			// if err != nil {
			// 	b.Error("Error", err)
			// }
		}
	})
}
