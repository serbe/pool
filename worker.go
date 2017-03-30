package pool

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		crawl(task)
	}
}
