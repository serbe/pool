package pool

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		p.dec()
		crawl(task)
		p.ResultChan <- task.result
		p.inc()
	}
}
