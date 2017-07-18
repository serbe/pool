package pool

import (
	"net/url"
	"sync"
	"time"
)

var (
	t10ms   = time.Duration(10) * time.Millisecond
	timeout = time.Duration(5) * time.Second
)

// Pool - pool of goroutines
type Pool struct {
	m            sync.RWMutex
	numWorkers   int
	freeWorkers  int
	finishedJobs int
	inputJobs    int
	startTime    time.Time
	workChan     chan *Task
	inputChan    chan *Task
	ResultChan   chan Task
	quitChan     chan bool
	endTaskChan  chan bool
	queue        taskList
}

// New - create new pool
func New(numWorkers int) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	p.freeWorkers = numWorkers
	p.finishedJobs = 0
	p.inputJobs = 0
	p.startTime = time.Now()
	p.workChan = make(chan *Task, numWorkers)
	p.inputChan = make(chan *Task)
	p.ResultChan = make(chan Task)
	p.endTaskChan = make(chan bool, numWorkers)
	p.quitChan = make(chan bool)
	for i := 0; i < numWorkers; i++ {
		go p.worker(i)
	}
	go p.run()
	return p
}

// Add - add new task to pool
func (p *Pool) Add(target string, proxy string) error {
	t := new(Task)
	targetURL, err := url.Parse(target)
	if err != nil {
		// log.Println("Error in Add Parse target", target, err)
		return err
	}
	t.Target = targetURL
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		// log.Println("Error in Add Parse proxy", proxy, err)
		return err
	}
	t.Proxy = proxyURL
	p.pushTask(t)
	return nil
}

func (p *Pool) run() {
runLoop:
	for {
		select {
		case work := <-p.inputChan:
			_ = p.queue.put(work)
			p.popTask()
			// if err != nil {
			// log.Println("Error in p.queue.put", err)
			// }
		case <-p.endTaskChan:
			p.popTask()
		// case <-time.After(t10ms):
		// 	p.popTask()
		case <-p.quitChan:
			close(p.ResultChan)
			break runLoop
			// if p.free() > 0 {
			// 	if p.queue.length() > 0 {
			// 		work, _ := p.queue.get()
			// 		// if err == nil {
			// 		p.workChan <- work
			// 		// } else {
			// 		// log.Println("Error in p.queue.get", err)
			// 		// }
			// 	} else if p.finishedJobs > 0 && p.finishedJobs == p.inputJobs {
			// 		if p.inputJobs == 1 && time.Since(p.startTime) > timeout || p.inputJobs != 1 {
			// 			close(p.ResultChan)
			// 			break runLoop
			// 		}
			// 	}
			// }
		}
	}
}

// SetTimeout - set http client timeout in second
func (p *Pool) SetTimeout(t int) {
	timeout = time.Duration(t) * time.Second
}
