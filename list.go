package pool

import (
	"sync"
	"sync/atomic"
)

type taskList struct {
	sync.RWMutex
	list []Task
	len  int64
}

func (t *taskList) put(task Task) {
	t.Lock()
	t.list = append(t.list, task)
	t.len++
	t.Unlock()
}

func (t *taskList) get() (Task, bool) {
	var task Task
	if t.length() > 0 {
		t.Lock()
		task = t.list[0]
		t.list = t.list[1:]
		t.len--
		t.Unlock()
		return task, true
	}
	return task, false
}

func (t *taskList) length() int64 {
	return atomic.LoadInt64(&t.len)
}
