package pool

import (
	"log"
	"testing"
	"time"
)

func Test_Pool(t *testing.T) {
	pool := New(2)
	if pool.numWorkers != 2 {
		t.Errorf("Expected %v, got %v", 2, pool.numWorkers)
	}
	err := pool.Add("", "1:")
	if err == nil {
		t.Errorf("Expected %v, got %v", err, nil)
	}
	err = pool.Add("1:", "")
	if err == nil {
		t.Errorf("Expected %v, got %v", err, nil)
	}
	err = pool.Add("https://ya.ru/", "")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	if pool.inputJobs != 1 {
		t.Errorf("Expected %v, got %v", 1, pool.inputJobs)
	}
	pool.SetHTTPTimeout(10)
	if timeout != time.Duration(10)*time.Second {
		t.Errorf("Expected %v, got %v", time.Duration(10)*time.Second, timeout)
	}
	err = pool.Add("https://ya.ru/", "http://bing.com/search?q=dotnet")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	if pool.inputJobs != 2 {
		t.Errorf("Expected %v, got %v", 2, pool.inputJobs)
	}

	pool.SetTaskTimeout(1)

	for result := range pool.ResultChan {
		log.Println(result.ID)
		log.Println(len(result.Body))
	}
}
