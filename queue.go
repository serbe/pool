package pool

import (
	"sync"
)

type ringQueue struct {
	sync.RWMutex
	nodes []Task
	head  int
	tail  int
	cnt   int
}

func newRingQueue() ringQueue {
	return ringQueue{
		nodes: make([]Task, 2),
	}
}

func (q *ringQueue) resize(n int) {
	nodes := make([]Task, n)
	if q.head < q.tail {
		copy(nodes, q.nodes[q.head:q.tail])
	} else {
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.tail])
	}

	q.tail = q.cnt % n
	q.head = 0
	q.nodes = nodes
}

func (q *ringQueue) put(task Task) {
	q.Lock()
	if q.cnt == len(q.nodes) {
		q.resize(q.cnt * 2)
	}
	q.nodes[q.tail] = task
	q.tail = (q.tail + 1) % len(q.nodes)
	q.cnt++
	q.Unlock()
}

func (q *ringQueue) get() (Task, bool) {
	var task Task
	q.Lock()
	if q.cnt == 0 {
		q.Unlock()
		return task, false
	}
	task = q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.cnt--

	if n := len(q.nodes) / 2; n > 2 && q.cnt <= n {
		q.resize(n)
	}
	q.Unlock()
	return task, true
}

func (q *ringQueue) Cap() int {
	return cap(q.nodes)
}

func (q *ringQueue) Len() int {
	return q.cnt
}
