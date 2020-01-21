package pool

import (
	"time"
)

// Task - structure describing a task
type Task struct {
	ID           int64
	Hostname     string
	Proxy        string
	Body         []byte
	ResponseTime time.Duration
	Error        error
}
