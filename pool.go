package pool

import (
	"log"
	"time"
)

var timeout int64

// Pool - specification of golang pool
type Pool struct {
	addedTasks int64
	in         ringQueue
	out        ringQueue
	toWork     chan Task
	fromWork   chan Task
	quit       chan struct{}
	workers    []worker
	// numWorkers     int64
	// useQuitTimeout bool
	// waitingTasks   uint32
	// runningPool    uint32
	// freeWorkers    int64
	// completedTasks int64
	// timeout        time.Duration
	// quitTimeout    time.Duration
	// timer          time.Timer
}

// New - create new goroutine pool with channels
// numWorkers - max workers
func New(numWorkers int64) *Pool {
	var (
		i       int64
		workers []worker
	)
	p := &Pool{
		in:       newRingQueue(),
		out:      newRingQueue(),
		toWork:   make(chan Task, numWorkers),
		fromWork: make(chan Task, numWorkers),
	}
	for i < numWorkers {
		worker := worker{
			id:   i,
			in:   p.toWork,
			out:  p.fromWork,
			quit: make(chan struct{}),
		}
		go worker.start()
		workers = append(workers, worker)
		i++
	}
	p.workers = workers
	go p.start()
	// p.freeWorkers = numWorkers
	// p.workChan = make(chan Task)
	// p.inputTaskChan = make(chan Task, 1)
	// p.ResultChan = make(chan Task, 1)
	// p.endTaskChan = make(chan struct{}, 1)
	// p.quit = make(chan struct{}, 1)
	// p.queue = newRingQueue()
	// p.timeout = time.Duration(10) * time.Second
	// go p.runBroker()
	// // go p.runWorkers()
	// p.runningPool = 1
	// p.waitingTasks = 1
	return p
}

func (p *Pool) start() {
	tick := time.Tick(time.Duration(200) * time.Microsecond)
	for {
		select {
		case <-tick:
			task, ok := p.in.get()
			if ok {
				p.toWork <- task
			}
		case task := <-p.fromWork:
			log.Println("pool get task", task.ID)
			p.out.put(task)
		}
	}
	// for {
	// 	select {
	// 	case task := <-p.inputTaskChan:
	// 		task.ID = p.GetAddedTasks()
	// 		p.addTask(task)
	// 	case <-p.endTaskChan:
	// 		// p.incWorkers()
	// 		p.tryGetTask()
	// 	case <-p.quit:
	// 		atomic.StoreUint32(&p.runningPool, 0)
	// 		p.EndWaitingTasks()
	// 		close(p.workChan)
	// 		close(p.ResultChan)
	// 		break
	// 	}
	// }
}

// Add - adding task to pool
func (p *Pool) Add(hostname string, proxy string) error {
	if hostname == "" {
		return errEmptyTarget
	}
	// if !p.poolIsRunning() {
	// 	return errNotRun
	// }
	// if !p.poolIsWaitingTasks() {
	// 	return errNotWait
	// }
	task := Task{
		Hostname: hostname,
		Proxy:    proxy,
	}
	p.in.put(task)
	return nil
}

// Get - try to getting finished task from pool
func (p *Pool) Get() (Task, bool) {
	return p.out.get()
}

// // Quit - send quit signal to pool
// func (p *Pool) Quit() {
// 	atomic.StoreUint32(&p.runningPool, 0)
// 	p.EndWaitingTasks()
// 	p.quit <- struct{}{}
// }

// func (p *Pool) poolIsRunning() bool {
// 	return atomic.LoadUint32(&p.runningPool) != 0
// }

// // EndWaitingTasks - stop pool waiting tasks
// func (p *Pool) EndWaitingTasks() {
// 	atomic.StoreUint32(&p.waitingTasks, 0)
// }

// func (p *Pool) poolIsWaitingTasks() bool {
// 	return atomic.LoadUint32(&p.waitingTasks) == 1
// }
