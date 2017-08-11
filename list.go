package pool

import (
	"errors"
	"sync"
)

var (
	// errEmptyTaskList = errors.New("task list is empty")
	errNilTask = errors.New("task is nil")
)

type taskList struct {
	m   sync.Mutex
	len int
	val []Task
}

func (t *taskList) put(task Task) {
	t.m.Lock()
	t.val = append(t.val, task)
	t.len++
	t.m.Unlock()
}

func (t *taskList) get() (Task, bool) {
	t.m.Lock()
	var task Task
	if t.len > 0 {
		task = t.val[0]
		t.len--
		t.val = t.val[1:]
		t.m.Unlock()
		return task, true
	}
	t.m.Unlock()
	return task, false
}
