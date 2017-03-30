package pool

import (
	"errors"
	"sync"
)

var (
	errEmptyTaskList = errors.New("task list is empty")
	errNilTask       = errors.New("task is nil")
)

type taskList struct {
	m   sync.RWMutex
	len int
	val []*task
}

func (t *taskList) put(task *task) error {
	t.m.Lock()
	defer t.m.Unlock()
	if task == nil {
		return errNilTask
	}
	t.val = append(t.val, task)
	t.len++
	return nil
}

func (t *taskList) get() (*task, error) {
	t.m.Lock()
	defer t.m.Unlock()
	var task *task
	if t.len > 0 {
		task = t.val[0]
		t.len--
		t.val = t.val[1:]
		return task, nil
	}
	return nil, errEmptyTaskList
}

func (t *taskList) length() int {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.len
}
