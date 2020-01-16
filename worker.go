package pool

import "sync/atomic"

type Worker struct {
	id   int64
	pool *Pool
	in   chan Task
	out  chan TaskResult
	quit chan struct{}
}

func (w *Worker) start() {
	go func() {
		for {
			select {
			case task := <-w.in:
				w.out <- w.crawl(task)
			case <-w.quit:
				break
			}
		}
	}()
}

// func (p *Pool) worker(id int64) {
// 	for task := range p.workChan {
// 		if p.useQuitTimeout {
// 			p.timer.Stop()
// 		}
// 		taskResult := p.crawl(task)
// 		if p.poolIsRunning() {
// 			p.ResultChan <- taskResult
// 			p.endTaskChan <- struct{}{}
// 			p.incCompletedTasks()
// 			if !p.poolIsWaitingTasks() && p.GetAddedTasks() == p.GetCompletedTasks() {
// 				atomic.StoreUint32(&p.runningPool, 0)
// 				p.quit <- struct{}{}
// 				break
// 			}
// 			if p.useQuitTimeout && p.isCompleteJobs() {
// 				p.timer.Reset(p.quitTimeout)
// 			}
// 		} else {
// 			break
// 		}
// 	}
// }

// func (p *Pool) runWorkers() {
// 	var i int64
// 	for i = 0; i < p.numWorkers; i++ {
// 		go p.worker(i)
// 	}
// }

// func (p *Pool) getFreeWorkers() int64 {
// 	return atomic.LoadInt64(&p.freeWorkers)
// }

func (p *Pool) incWorkers() {
	atomic.AddInt64(&p.freeWorkers, 1)
}

func (p *Pool) decWorkers() {
	atomic.AddInt64(&p.freeWorkers, -1)
}
