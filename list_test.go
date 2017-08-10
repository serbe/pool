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
	_, ok := list.get()
	if !ok {
		t.Errorf("Expected %v, got %v", true, err)
	}
	_, ok = list.get()
	if ok {
		t.Errorf("Expected %v, got %v", false, err)
	}
	length = list.length()
	if length != 0 {
		t.Errorf("Expected %v, got %v", 0, length)
	}
}
