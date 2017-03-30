package pool

// Pool - pool of goroutines
type Pool struct {
	numWorkers int
	workChan   chan *task
	// tasks      []*task
}

// New - create new pool
func New(numWorkers int) *Pool {
	p := new(Pool)
	p.numWorkers = numWorkers
	for i := 0; i < numWorkers; i++ {
		go p.worker(i)
	}
	return p
}
