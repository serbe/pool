package pool

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

var (
	errEmptyTarget = errors.New("error: empty target hostname")
	errNotRun      = errors.New("error: pool is not running")
	errNotWait     = errors.New("error: pool is not waiting tasks")
)

// Task - structure describing a task
type Task struct {
	ID           int64
	WorkerID     int64
	Hostname     string
	Body         []byte
	Proxy        *url.URL
	Response     *http.Response
	ResponceTime time.Duration
	Error        error
}

// Add - add new task to pool
func (p *Pool) Add(hostname string, proxy *url.URL) error {
	if hostname == "" {
		return errEmptyTarget
	}
	if !p.poolIsRunning() {
		return errNotRun
	}
	if !p.poolIsWaitingTasks() {
		return errNotWait
	}
	_, err := url.Parse(hostname)
	if err != nil {
		return err
	}
	task := &Task{
		Hostname: hostname,
		Proxy:    proxy,
	}
	p.inputTaskChan <- task
	return nil
}

func (p *Pool) addTask(task *Task) {
	if p.getFreeWorkers() > 0 {
		p.decWorkers()
		p.workChan <- task
	} else {
		p.queue.put(task)
	}
}

func (p *Pool) tryGetTask() {
	task, ok := p.queue.get()
	if ok {
		p.decWorkers()
		p.workChan <- task
	}
}

// SetTimeout - set http timeout in millisecond
func (p *Pool) SetTimeout(t int64) {
	p.timeout = time.Duration(t) * time.Millisecond
}

// SetQuitTimeout - set timeout to quit after finish all tasks in millisecond
func (p *Pool) SetQuitTimeout(t int64) {
	p.useQuitTimeout = true
	p.quitTimeout = time.Duration(t) * time.Millisecond
	p.timer = time.NewTimer(p.quitTimeout)
	go func() {
		<-p.timer.C
		p.quit <- true
	}()
}

func (p *Pool) crawl(t *Task) *Task {
	startTime := time.Now()
	client := &http.Client{
		Timeout: p.timeout,
	}
	if t.Proxy != nil {
		client.Transport = &http.Transport{
			Proxy:             http.ProxyURL(t.Proxy),
			DisableKeepAlives: true,
		}
	}
	req, err := http.NewRequest("GET", t.Hostname, nil)
	if err != nil {
		t.Error = err
		return t
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		t.Error = err
		return t
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error = err
		return t
	}
	t.Body = body
	t.Response = resp
	t.ResponceTime = time.Since(startTime)
	return t
}

// GetAddedTasks - get num of added tasks
func (p *Pool) GetAddedTasks() int64 {
	return atomic.LoadInt64(&p.addedTasks)
}

func (p *Pool) incAddedTasks() {
	atomic.AddInt64(&p.addedTasks, 1)
}

// GetCompletedTasks - get num of completed tasks
func (p *Pool) GetCompletedTasks() int64 {
	return atomic.LoadInt64(&p.completedTasks)
}

func (p *Pool) incCompletedTasks() {
	atomic.AddInt64(&p.completedTasks, 1)
}

func (p *Pool) isCompleteJobs() bool {
	return p.GetCompletedTasks() > 0 && p.GetCompletedTasks() == p.GetAddedTasks()
}
