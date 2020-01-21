package pool

import (
	"sync"
	"time"
)

var timeout int64

// Pool - specification of golang pool
type Pool struct {
	running        bool
	useOutChan     bool
	numWorkers     int64
	addedTasks     int64
	completedTasks int64
	in             ringQueue
	out            ringQueue
	toWork         chan Task
	fromWork       chan Task
	outTasks       chan Task
	quit           chan struct{}
	workers        []worker
	wg             sync.WaitGroup
	taskWG         sync.WaitGroup
}

// New - create new goroutine pool with channels
// numWorkers - max workers
func New(numWorkers int64) *Pool {
	// rand.Seed(time.Now().UnixNano())
	p := &Pool{
		numWorkers: numWorkers,
		in:         newRingQueue(),
		out:        newRingQueue(),
		toWork:     make(chan Task, numWorkers),
		fromWork:   make(chan Task, numWorkers),
		outTasks:   make(chan Task, 1),
		quit:       make(chan struct{}),
	}
	p.wg.Add(1)
	p.startWorkers()
	go p.start()
	p.wg.Wait()
	return p
}

func (p *Pool) startWorkers() {
	var (
		i       int64
		workers []worker
	)
	for i < p.numWorkers {
		p.wg.Add(1)
		worker := worker{
			id:   i,
			in:   p.toWork,
			out:  p.fromWork,
			quit: make(chan struct{}),
			wg:   &p.wg,
		}
		go worker.start()
		workers = append(workers, worker)
		i++
	}
	p.workers = workers
}

func (p *Pool) start() {
	p.running = true
	p.wg.Done()
	tick := time.Tick(time.Duration(200) * time.Microsecond)
	for {
		select {
		case <-tick:
			task, ok := p.in.get()
			if ok {
				p.toWork <- task
			}
		case task := <-p.fromWork:
			if !p.useOutChan {
				p.out.put(task)
			} else {
				p.outTasks <- task
			}
			p.completedTasks++
			p.taskWG.Done()
		case <-p.quit:
			for i := range p.workers {
				p.wg.Add(1)
				p.workers[i].stop()
			}
			// close(p.quit)
			break
		}
	}
}

// Add - adding task to pool
func (p *Pool) Add(hostname string, proxy string) error {
	if hostname == "" {
		return errEmptyTarget
	}
	if !p.running {
		return errNotRun
	}
	task := Task{
		ID:       p.addedTasks,
		Hostname: hostname,
		Proxy:    proxy,
	}
	p.addedTasks++
	p.in.put(task)
	p.taskWG.Add(1)
	return nil
}

// Get - try to getting finished task from pool
func (p *Pool) Get() (Task, bool) {
	return p.out.get()
}

// Stop - stop pool and all workers
func (p *Pool) Stop() {
	p.quit <- struct{}{}
	p.wg.Wait()
	p.running = false
}

// NetTimeout - set crawl timeout in milliseconds
func (p *Pool) NetTimeout(millis int64) {
	timeout = millis
}

// IsRunning - check pool status is running
func (p *Pool) IsRunning() bool {
	return p.running
}

// Wait - wait all task is done
func (p *Pool) Wait() {
	p.taskWG.Wait()
}

// Added - number of adding tasks
func (p *Pool) Added() int64 {
	return p.addedTasks
}

// Completed - number of completed tasks
func (p *Pool) Completed() int64 {
	return p.completedTasks
}

// UseOutChan - use chan to get results
func (p *Pool) UseOutChan() chan Task {
	p.useOutChan = true
	return p.outTasks
}

// IsUseOutChan - get status of use chan to get results
func (p *Pool) IsUseOutChan() bool {
	return p.useOutChan
}
