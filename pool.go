package pool

// Pool - pool of goroutines
type Pool struct {
	numWorkers int
	workChan   chan *task
	queue      taskList
}

// New - create new pool
func New(numWorkers int) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	p.workChan = make(chan *task, numWorkers)
	for i := 0; i < numWorkers; i++ {
		go p.worker(i)
	}
	return p
}

// Add - add new task to pool
func (p *Pool) Add(address string) {
	t := new(task)
	p.queue.put(t)
}

func (p *Pool) run() {
	for {

	}
}
