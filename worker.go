package pool

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		if p.timerIsRunning {
			p.timer.Stop()
		}
		p.dec()
		task.WorkerID = id
		crawl(task)
		p.ResultChan <- *task
		p.inc()
		p.endTaskChan <- true
		if p.timerIsRunning && p.freeWorkers == p.numWorkers {
			p.timer.Reset(p.quitTimeout)
		}
	}
}
