package pool

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		p.dec()
		crawl(task)
		p.ResultChan <- Result{Body: task.result, Address: task.address}
		p.inc()
	}
}
