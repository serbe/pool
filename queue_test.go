package pool

import (
	"testing"
)

func TestQueue(t *testing.T) {
	queue := newRingQueue()
	queue.put(Task{Hostname: "1"})
	if queue.Len() != 1 {
		t.Errorf("Got %v queue len, want %v", queue.Len(), 1)
	}
	if queue.Cap() != 2 {
		t.Errorf("Got %v queue cap, want %v", queue.Cap(), 2)
	}
	queue.put(Task{Hostname: "2"})
	if queue.Len() != 2 {
		t.Errorf("Got %v queue len, want %v", queue.Len(), 2)
	}
	if queue.Cap() != 2 {
		t.Errorf("Got %v queue cap, want %v", queue.Cap(), 2)
	}
	queue.put(Task{Hostname: "3"})
	if queue.Len() != 3 {
		t.Errorf("Got %v queue len, want %v", queue.Len(), 3)
	}
	if queue.Cap() != 4 {
		t.Errorf("Got %v queue cap, want %v", queue.Cap(), 4)
	}
	queue.put(Task{Hostname: "4"})
	queue.put(Task{Hostname: "5"})
	queue.put(Task{Hostname: "6"})
	queue.put(Task{Hostname: "7"})
	queue.put(Task{Hostname: "8"})
	queue.put(Task{Hostname: "9"})
	if queue.Len() != 9 {
		t.Errorf("Got %v queue len, want %v", queue.Len(), 9)
	}
	if queue.Cap() != 16 {
		t.Errorf("Got %v queue cap, want %v", queue.Cap(), 16)
	}
	task, ok := queue.get()
	if !ok {
		t.Errorf("Got %v in queue get, want %v", ok, true)
	}
	if task.Hostname != "1" {
		t.Errorf("Got %v task hostname, want %v", task.Hostname, "1")
	}
	task, ok = queue.get()
	if !ok {
		t.Errorf("Got %v in queue get, want %v", ok, true)
	}
	if task.Hostname != "2" {
		t.Errorf("Got %v task hostname, want %v", task.Hostname, "2")
	}
	task, ok = queue.get()
	if !ok {
		t.Errorf("Got %v in queue get, want %v", ok, true)
	}
	if task.Hostname != "3" {
		t.Errorf("Got %v task hostname, want %v", task.Hostname, "3")
	}
	_, _ = queue.get()
	_, _ = queue.get()
	_, _ = queue.get()
	_, _ = queue.get()
	_, _ = queue.get()
	task, _ = queue.get()
	if task.Hostname != "9" {
		t.Errorf("Got %v task hostname, want %v", task.Hostname, "9")
	}
	_, ok = queue.get()
	if ok {
		t.Errorf("Got %v in queue get, want %v", ok, false)
	}
}

func BenchmarkQueue(b *testing.B) {
	queue := newRingQueue()
	b.ResetTimer()

	n := b.N
	for i := 0; i < n; i++ {
		testTask := Task{
			Hostname: "",
		}
		queue.put(testTask)
	}
	for i := 0; i < n; i++ {
		_, ok := queue.get()
		if !ok {
			b.Errorf("Got %v error, want %v", ok, true)
		}
	}
}

func BenchmarkParallelQueue(b *testing.B) {
	queue := newRingQueue()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testTask := Task{
				Hostname: "",
			}
			queue.put(testTask)
		}
	})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, ok := queue.get()
			if !ok {
				b.Errorf("Got %v error, want %v", ok, true)
			}
		}
	})
}
