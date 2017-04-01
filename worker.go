package pool

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		p.dec()
		task.ID = id
		crawl(task)
		p.ResultChan <- *task
		p.inc()
	}
}
