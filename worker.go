package pool

import (
	"sync"
)

type worker struct {
	id      int64
	running bool
	in      chan Task
	out     chan Task
	quit    chan struct{}
	wg      *sync.WaitGroup
}

func (w *worker) start() {
	w.running = true
	w.wg.Done()
	for {
		select {
		case task := <-w.in:
			w.out <- crawl(task)
		case <-w.quit:
			// close(w.quit)
			w.running = false
			w.wg.Done()
			return
		}
	}
}

func (w *worker) stop() {
	w.quit <- struct{}{}
}
