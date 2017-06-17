package pool

import (
	"log"
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
	for i := 0; i < numWorkers; i++ {
		go p.worker(i)
	}
	go p.run()
	return p
}

// Add - add new task to pool
func (p *Pool) Add(target string, proxy string) {
	t := new(Task)
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Println("Error in Add Parse target", target, err)
		return
	}
	t.Target = targetURL
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Println("Error in Add Parse proxy", proxy, err)
		return
	}
	t.Proxy = proxyURL
	p.inputJobs++
	p.inputChan <- t
}

func (p *Pool) run() {
runLoop:
	for {
		select {
		case work := <-p.inputChan:
			err := p.queue.put(work)
			if err != nil {
				log.Println("Error in p.queue.put", err)
			}
		case <-time.After(t10ms):
			if p.free() > 0 {
				if p.queue.length() > 0 {
					work, err := p.queue.get()
					if err == nil {
						p.workChan <- work
					} else {
						log.Println("Error in p.queue.get", err)
					}
				} else if p.finishedJobs > 0 && p.finishedJobs == p.inputJobs {
					if p.inputJobs == 1 && time.Since(p.startTime) > timeout || p.inputJobs != 1 {
						close(p.ResultChan)
						break runLoop
					}
				}
			}
		}
	}
}

func (p *Pool) free() int {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.freeWorkers
}

func (p *Pool) inc() {
	p.m.Lock()
	p.freeWorkers++
	p.finishedJobs++
	p.m.Unlock()
}

func (p *Pool) dec() {
	p.m.Lock()
	p.freeWorkers--
	p.m.Unlock()
}

// SetTimeout - set http client timeout in second
func (p *Pool) SetTimeout(t int) {
	timeout = time.Duration(t) * time.Second
}
