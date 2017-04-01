package pool

import (
	"time"
)

// Task - structure describing a task
type Task struct {
	ID           int
	Address      string
	Proxy        string
	Body         []byte
	ResponceTime time.Duration
	Error        error
}
