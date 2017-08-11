package pool

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		task.WorkerID = id
		task = crawl(task)
		p.ResultChan <- task
		p.endTaskChan <- true
	}
}

func (p *Pool) runWorkers() {
	for i := 0; i < p.numWorkers; i++ {
		go p.worker(i)
	}
}
