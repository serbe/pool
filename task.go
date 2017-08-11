package pool

import (
	"net/http"
	"net/url"
	"time"
)

// Task - structure describing a task
type Task struct {
	ID           int
	WorkerID     int
	Hostname     string
	Proxy        *url.URL
	Response     *http.Response
	Body         []byte
	ResponceTime time.Duration
	Error        error
}

// func (p *Pool) popTask() {
// 	if p.freeWorkers > 0 {
// 		task, ok := p.queue.get()
// 		log.Println("poptask sucess p.queue.get", task.ID, task.Hostname)
// 		if ok {
// 			log.Println("poptask ok, try p.workChan <- task", task.ID, task.Hostname)
// 			p.workChan <- task
// 			log.Println("poptask sucess p.workChan <- task")
// 		}
// 	}
// }

func (p *Pool) inc() {
	p.m.Lock()
	p.freeWorkers++
	p.m.Unlock()
}

func (p *Pool) dec() {
	p.m.Lock()
	p.freeWorkers--
	p.m.Unlock()
}
