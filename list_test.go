package pool

import (
	"testing"
)

func Test_List(t *testing.T) {
	var list = new(taskList)
	var task = new(Task)
	err := list.put(task)
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	err = list.put(nil)
	if err != errNilTask {
		t.Errorf("Expected %v, got %v", errNilTask, err)
	}
	length := list.length()
	if length != 1 {
		t.Errorf("Expected %v, got %v", 1, length)
	}
	_, err = list.get()
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	_, err = list.get()
	if err != errEmptyTaskList {
		t.Errorf("Expected %v, got %v", errEmptyTaskList, err)
	}
	length = list.length()
	if length != 0 {
		t.Errorf("Expected %v, got %v", 0, length)
	}
}
