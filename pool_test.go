package pool

import (
	"net/url"
	"testing"
	"time"
)

func Test_Pool(t *testing.T) {
	pool := New(2)
	if pool.numWorkers != 2 {
		t.Errorf("Expected %v, got %v", 2, pool.numWorkers)
	}
	testURL := new(url.URL)
	err := pool.Add("", testURL)
	if err == nil {
		t.Errorf("Expected %v, got %v", err, nil)
	}
	err = pool.Add("1:", testURL)
	if err == nil {
		t.Errorf("Expected %v, got %v", err, nil)
	}
	err = pool.Add("https://ya.ru/", testURL)
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	if pool.getJobs() != 1 {
		t.Errorf("Expected %v, got %v", 1, pool.getJobs())
	}
	err = pool.Add("https://ya.ru/", testURL)
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	pool.SetHTTPTimeout(10)
	if timeout != time.Duration(10)*time.Second {
		t.Errorf("Expected %v, got %v", time.Duration(10)*time.Second, timeout)
	}
	testURL, err = url.Parse("http://bing.com/search?q=dotnet")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	err = pool.Add("https://ya.ru/", testURL)
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
	pool.SetTaskTimeout(1)
	if pool.getJobs() != 3 {
		t.Errorf("Expected %v, got %v", 3, pool.getJobs())
	}
	result := <-pool.ResultChan
	if result.Error != nil {
		t.Errorf("Expected %v, got %v", nil, result.Error)
	}
	<-pool.ResultChan
	<-pool.ResultChan
	pool.Quit()
}
