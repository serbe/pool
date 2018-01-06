package pool

import (
	"sync/atomic"
)

func (p *Pool) worker(id int64) {
	for task := range p.workChan {
		task.WorkerID = id
		task = crawl(task)
		p.ResultChan <- task
		p.endTaskChan <- true
	}
}

func (p *Pool) runWorkers() {
	var i int64
	for i = 0; i < p.numWorkers; i++ {
		go p.worker(i)
	}
}

func (p *Pool) getFreeWorkers() int64 {
	return atomic.LoadInt64(&p.freeWorkers)
}

func (p *Pool) incWorkers() {
	atomic.AddInt64(&p.freeWorkers, 1)
}

func (p *Pool) decWorkers() {
	atomic.AddInt64(&p.freeWorkers, -1)
}
