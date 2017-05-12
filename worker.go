package pool

import "log"

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		p.dec()
		task.ID = id
		err := crawl(task)
		if err != nil {
			log.Println("Error in crawl", err)
		}
		p.ResultChan <- *task
		p.inc()
	}
}
