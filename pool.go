package pool

import (
	"log"
	"net/url"
	"sync"
	"time"
)

var timeout = time.Duration(5) * time.Second

// Pool - pool of goroutines
type Pool struct {
	m              sync.RWMutex
	timerIsRunning bool
	numWorkers     int
	freeWorkers    int
	inputJobs      int
	quitTimeout    time.Duration
	workChan       chan *Task
	inputChan      chan *Task
	ResultChan     chan Task
	quitChan       chan bool
	endTaskChan    chan bool
	queue          taskList
	timer          *time.Timer
}

// New - create new pool
func New(numWorkers int) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	p.freeWorkers = numWorkers
	p.inputJobs = 0
	p.workChan = make(chan *Task, numWorkers+1)
	p.inputChan = make(chan *Task)
	p.ResultChan = make(chan Task)
	p.endTaskChan = make(chan bool, numWorkers+1)
	p.quitChan = make(chan bool)
	for i := 0; i < numWorkers; i++ {
		go p.worker(i)
	}
	go p.run()
	return p
}

// Add - add new task to pool
func (p *Pool) Add(hostname string, proxy *url.URL) {
	t := new(Task)
	t.Hostname = hostname
	t.Proxy = proxy
	p.inputJobs++
	t.ID = p.inputJobs
	log.Println("try to pushtask", t.Hostname)
	p.inputChan <- t
	log.Println("sucess pushtask", t.Hostname)
}

func (p *Pool) run() {
runLoop:
	for {
		select {
		case work := <-p.inputChan:
			_ = p.queue.put(work)
			p.popTask()
		case <-p.endTaskChan:
			p.popTask()
		case <-p.quitChan:
			close(p.ResultChan)
			close(p.workChan)
			break runLoop
		}
	}
}

// SetHTTPTimeout - set http client timeout in second
func (p *Pool) SetHTTPTimeout(t int) {
	timeout = time.Duration(t) * time.Second
}

// SetTaskTimeout - set task timeout in second before send quit signal
func (p *Pool) SetTaskTimeout(t int) {
	p.quitTimeout = time.Duration(t) * time.Second
	p.timer = time.NewTimer(p.quitTimeout)
	p.timerIsRunning = true
	go func() {
		<-p.timer.C
		p.quitChan <- true
		log.Println("End with timeout")
	}()
}

// Quit - send quit signal to pool
func (p *Pool) Quit() {
	p.quitChan <- true
}
