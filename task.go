package pool

import (
	"net/http"
	"net/url"
	"time"
)

// Task - structure describing a task
type Task struct {
	ID           int
	Target       *url.URL
	Proxy        *url.URL
	Response     *http.Response
	Body         []byte
	ResponceTime time.Duration
	Error        error
}

func (p *Pool) pushTask(t *Task) {
	p.inputJobs++
	p.inputChan <- t
}

func (p *Pool) popTask() {
	if p.freeWorkers > 0 && p.queue.length() > 0 {
		work, _ := p.queue.get()
		p.workChan <- work
	}
}

// func (p *Pool) free() int {
// p.m.RLock()
// defer p.m.RUnlock()
// return p.freeWorkers
// }

func (p *Pool) inc() {
	p.m.Lock()
	p.freeWorkers++
	p.finishedJobs++
	p.m.Unlock()
}

func (p *Pool) dec() {
	p.m.Lock()
	p.freeWorkers--
	p.m.Unlock()
}
