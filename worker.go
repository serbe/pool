package pool

import (
	"log"
)

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		log.Println("Start task ", id)
		p.dec()
		task.ID = id
		crawl(task)
		p.ResultChan <- *task
		p.inc()
		p.endTaskChan <- true
		log.Println("Finish task ", id)
	}
}
