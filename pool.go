package pool

import (
	"sync/atomic"
	"time"
)

// Pool - specification of gopool
type Pool struct {
	useQuitTimeout bool
	waitingTasks   uint32
	runningPool    uint32
	numWorkers     int64
	freeWorkers    int64
	addedTasks     int64
	completedTasks int64
	quit           chan bool
	endTaskChan    chan bool
	workChan       chan *Task
	inputTaskChan  chan *Task
	ResultChan     chan *Task
	queue          *taskQueue
	timeout        time.Duration
	quitTimeout    time.Duration
	timer          *time.Timer
}

// New - create new goroutine pool with channels
// numWorkers - max workers
func New(numWorkers int64) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	p.freeWorkers = numWorkers
	p.workChan = make(chan *Task)
	p.inputTaskChan = make(chan *Task, 1)
	p.ResultChan = make(chan *Task, 1)
	p.endTaskChan = make(chan bool, 1)
	p.quit = make(chan bool, 1)
	p.queue = new(taskQueue)
	p.timeout = time.Duration(10) * time.Second
	go p.runBroker()
	go p.runWorkers()
	p.runningPool = 1
	p.waitingTasks = 1
	return p
}

func (p *Pool) runBroker() {
loopPool:
	for {
		select {
		case task := <-p.inputTaskChan:
			p.incAddedTasks()
			task.ID = p.GetAddedTasks()
			p.addTask(task)
		case <-p.endTaskChan:
			p.incWorkers()
			p.tryGetTask()
		case <-p.quit:
			close(p.workChan)
			close(p.ResultChan)
			break loopPool
		}
	}
}

// Quit - send quit signal to pool
func (p *Pool) Quit() {
	atomic.StoreUint32(&p.runningPool, 0)
	p.quit <- true
	p.EndWaitingTasks()
}

func (p *Pool) poolIsRunning() bool {
	return atomic.LoadUint32(&p.runningPool) != 0
}

// EndWaitingTasks - set end pool waiting tasks
func (p *Pool) EndWaitingTasks() {
	atomic.StoreUint32(&p.waitingTasks, 0)
}

func (p *Pool) poolIsWaitingTasks() bool {
	return atomic.LoadUint32(&p.waitingTasks) == 1
}
