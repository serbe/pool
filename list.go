package pool

import (
	"sync"
)

type taskList struct {
	sync.RWMutex
	list []Task
}

func (t *taskList) put(task Task) {
	t.Lock()
	t.list = append(t.list, task)
	t.Unlock()
}

func (t *taskList) get() (Task, bool) {
	t.Lock()
	var task Task
	if len(t.list) > 0 {
		task = t.list[0]
		t.list = t.list[1:]
		t.Unlock()
		return task, true
	}
	t.Unlock()
	return task, false
}

func (t *taskList) length() int {
	t.RLock()
	length := len(t.list)
	t.RUnlock()
	return length
}
