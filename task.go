package pool

import (
	"net/http"
	"net/url"
	"time"
)

// Task - structure describing a task
type Task struct {
	ID           int
	WorkerID     int
	Hostname     string
	Proxy        *url.URL
	Response     *http.Response
	Body         []byte
	ResponceTime time.Duration
	Error        error
}
