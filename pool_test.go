package pool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_Pool(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test")
	}))
	defer ts.Close()

	pool := New(2)
	if pool.numWorkers != 2 {
		t.Errorf("Expected %v, got %v", 2, pool.numWorkers)
	}
	if pool.free() != 2 {
		t.Errorf("Expected %v, got %v", 2, pool.free())
	}
	err := pool.Add("", "1:")
	if err == nil {
		t.Errorf("Expected %v, got %v", err, nil)
	}
	err = pool.Add("1:", "")
	if err == nil {
		t.Errorf("Expected %v, got %v", err, nil)
	}
	err = pool.Add("/", "")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	if pool.inputJobs != 1 {
		t.Errorf("Expected %v, got %v", 1, pool.inputJobs)
	}
	pool.SetTimeout(10)
	if timeout != time.Duration(10)*time.Second {
		t.Errorf("Expected %v, got %v", time.Duration(10)*time.Second, timeout)
	}
	err = pool.Add("/", "127.0.0.1")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	if pool.inputJobs != 2 {
		t.Errorf("Expected %v, got %v", 2, pool.inputJobs)
	}

	// for result := range pool.ResultChan {
	// 	log.Println(result.ID)
	// 	log.Println(result.Body)
	// }
}
