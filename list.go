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
	m   sync.RWMutex
	len int
	val []*Task
}

func (t *taskList) put(task *Task) error {
	t.m.Lock()

	if task == nil {
		t.m.Unlock()
		return errNilTask
	}
	t.val = append(t.val, task)
	t.len++
	t.m.Unlock()
	return nil
}

func (t *taskList) get() (*Task, bool) {
	t.m.Lock()
	var task *Task
	if t.len > 0 {
		task = t.val[0]
		t.len--
		t.val = t.val[1:]
		t.m.Unlock()
		return task, true
	}
	t.m.Unlock()
	return nil, false
}

func (t *taskList) length() int {
	t.m.RLock()
	len := t.len
	t.m.RUnlock()
	return len
}
