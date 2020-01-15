package pool

import (
	"log"
	"sync/atomic"
	"time"
)

// Pool - specification of golang pool
type Pool struct {
	useQuitTimeout bool
	waitingTasks   uint32
	runningPool    uint32
	numWorkers     int64
	freeWorkers    int64
	addedTasks     int64
	completedTasks int64
	quit           chan struct{}
	waitingWorkers chan struct{}
	toWorker       chan Task
	fromWorker     chan TaskResult
	ResultChan     chan TaskResult
	workers        []*Worker
	queue          ringQueue
	timeout        time.Duration
	quitTimeout    time.Duration
	timer          time.Timer
}

// New - create new goroutine pool with channels
// numWorkers - max workers
func New(numWorkers int64) *Pool {
	p := &Pool{
		numWorkers:     numWorkers,
		freeWorkers:    numWorkers,
		toWorker:       make(chan Task, 1),
		fromWorker:     make(chan TaskResult, 1),
		ResultChan:     make(chan TaskResult, 1),
		workers:        make([]*Worker, 4),
		quit:           make(chan struct{}, 1),
		waitingWorkers: make(chan struct{}, 4),
		queue:          newRingQueue(),
		timeout:        time.Duration(10) * time.Second,
	}
	var i int64
	for i = 0; i < numWorkers; i++ {
		worker := &Worker{
			id:   i,
			pool: p,
			in:   p.toWorker,
			out:  p.fromWorker,
			quit: make(chan struct{}, 1),
		}
		p.workers[i] = worker
	}
	// go p.runBroker()
	// go p.runWorkers()
	// p.runningPool = 1
	p.waitingTasks = 1
	return p
}

func (p *Pool) Start() {
	atomic.StoreUint32(&p.runningPool, 1)
	// tick := time.Tick(100 * time.Millisecond)
	for {
		select {
		// case task := <-p.inputTaskChan:
		// 	task.ID = p.GetAddedTasks()
		// 	p.addTask(task)
		case result := <-p.fromWorker:
			p.ResultChan <- result
		case <-p.waitingWorkers:
			task, ok := p.queue.get()
			if !ok {
				log.Println("queue is empty")
				break
			}
			p.toWorker <- task
		case <-p.quit:
			atomic.StoreUint32(&p.runningPool, 0)
			p.EndWaitingTasks()
			var i int64
			for i = 0; i < p.numWorkers; i++ {
				p.workers[i].quit <- struct{}{}
			}
			close(p.ResultChan)
			break
		}
	}
}

// Add - add new task to pool
func (p *Pool) Add(hostname string, proxy string) error {
	if hostname == "" {
		return errEmptyTarget
	}
	// if !p.poolIsRunning() {
	// 	return errNotRun
	// }
	if !p.poolIsWaitingTasks() {
		return errNotWait
	}
	task := Task{
		Hostname: hostname,
		Proxy:    proxy,
	}
	p.incAddedTasks()
	p.queue.put(task)
	return nil
}

// Quit - send quit signal to pool
func (p *Pool) Quit() {
	atomic.StoreUint32(&p.runningPool, 0)
	p.EndWaitingTasks()
	p.quit <- struct{}{}
}

func (p *Pool) poolIsRunning() bool {
	return atomic.LoadUint32(&p.runningPool) != 0
}

// EndWaitingTasks - stop pool waiting tasks
func (p *Pool) EndWaitingTasks() {
	atomic.StoreUint32(&p.waitingTasks, 0)
}

func (p *Pool) poolIsWaitingTasks() bool {
	return atomic.LoadUint32(&p.waitingTasks) == 1
}
