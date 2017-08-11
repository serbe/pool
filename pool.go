package pool

import (
	"net/url"
	"time"
)

var (
	timeout = time.Duration(5) * time.Second
	t50ms   = time.Duration(50) * time.Millisecond
)

// Pool - pool of goroutines
type Pool struct {
	timerIsRunning bool
	numWorkers     int
	freeWorkers    int
	inputJobs      int
	quitTimeout    time.Duration
	workChan       chan Task
	inputTaskChan  chan Task
	ResultChan     chan Task
	quit           chan bool
	endTaskChan    chan bool
	queue          taskList
	timer          *time.Timer
}

// New - create new pool
func New(numWorkers int) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	p.freeWorkers = numWorkers
	p.workChan = make(chan Task)
	p.inputTaskChan = make(chan Task)
	p.ResultChan = make(chan Task)
	p.endTaskChan = make(chan bool)
	p.quit = make(chan bool)
	go p.run()
	for i := 0; i < numWorkers; i++ {
		go p.worker(i)
	}
	return p
}

// Add - add new task to pool
func (p *Pool) Add(hostname string, proxy *url.URL) error {
	if hostname == "" {
		return errNilTask
	}
	task := Task{
		Hostname: hostname,
		Proxy:    proxy,
	}
	p.inputTaskChan <- task
	return nil
}

func (p *Pool) run() {
loopPool:
	for {
		select {
		case task := <-p.inputTaskChan:
			p.inputJobs++
			task.ID = p.inputJobs
			if p.freeWorkers > 0 {
				if p.timerIsRunning {
					p.timer.Stop()
				}
				p.freeWorkers--
				p.workChan <- task
			} else {
				p.queue.put(task)
			}
		case <-p.endTaskChan:
			p.freeWorkers++
			if p.timerIsRunning && p.freeWorkers == p.numWorkers {
				p.timer.Reset(p.quitTimeout)
			}
		case <-p.quit:
			close(p.ResultChan)
			close(p.workChan)
			break loopPool
		case <-time.After(t50ms):
			if p.freeWorkers > 0 {
				task, ok := p.queue.get()
				if ok {
					if p.timerIsRunning {
						p.timer.Stop()
					}
					p.freeWorkers--
					p.workChan <- task
				}
			}
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
		p.quit <- true
		// log.Println("End with timeout")
	}()
}

// Quit - send quit signal to pool
func (p *Pool) Quit() {
	p.quit <- true
}
