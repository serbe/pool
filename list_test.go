package pool

import (
	"testing"
)

func Test_List(t *testing.T) {
	var list = new(taskList)
	var task Task
	list.put(task)
	length := list.length()
	if length != 1 {
		t.Errorf("Expected %v, got %v", 1, length)
	}
	_, ok := list.get()
	if !ok {
		t.Errorf("Expected %v, got %v", true, ok)
	}
	_, ok = list.get()
	if ok {
		t.Errorf("Expected %v, got %v", false, ok)
	}
	length = list.length()
	if length != 0 {
		t.Errorf("Expected %v, got %v", 0, length)
	}
}
