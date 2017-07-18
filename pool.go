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
	quitTimeout  time.Duration
	startTime    time.Time
	workChan     chan *Task
	inputChan    chan *Task
	ResultChan   chan Task
	quitChan     chan bool
	endTaskChan  chan bool
	queue        taskList
	timer        *time.Timer
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
		return err
	}
	t.Target = targetURL
	proxyURL, err := url.Parse(proxy)
	if err != nil {
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
		case <-p.endTaskChan:
			p.popTask()
		case <-p.quitChan:
			close(p.ResultChan)
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
	go func() {
		<-p.timer.C
		p.quitChan <- true
		log.Println("End with timeout")
	}()
}
