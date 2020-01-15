package pool

import (
	"testing"
)

func TestQueue(t *testing.T) {
	queue := newRingQueue()
	testTask := Task{
		Hostname: "",
	}
	queue.put(testTask)
	if queue.cnt != 1 {
		t.Errorf("Got %v error, want %v", queue.cnt, 1)
	}
	queue.put(testTask)
	_, ok := queue.get()
	if !ok {
		t.Errorf("Got %v error, want %v", ok, true)
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
