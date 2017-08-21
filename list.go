package pool

import (
	"sync"
)

type taskList struct {
	m    sync.RWMutex
	list []Task
}

func (t *taskList) put(task Task) {
	t.m.Lock()
	t.list = append(t.list, task)
	t.m.Unlock()
}

func (t *taskList) get() (Task, bool) {
	t.m.Lock()
	var task Task
	if len(t.list) > 0 {
		task = t.list[0]
		t.list = t.list[1:]
		t.m.Unlock()
		return task, true
	}
	t.m.Unlock()
	return task, false
}

func (t *taskList) length() int {
	t.m.RLock()
	length := len(t.list)
	t.m.RUnlock()
	return length
}
