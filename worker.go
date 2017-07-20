package pool

import "time"

func (p *Pool) worker(id int) {
	for task := range p.workChan {
		p.dec()
		task.WorkerID = id
		crawl(task)
		p.ResultChan <- *task
		p.inc()
		p.endTaskChan <- true
		if p.timerIsRunning {
			p.timer = time.NewTimer(p.quitTimeout)
		}
	}
}
