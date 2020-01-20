package pool

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func crawl(task Task) Task {
	startTime := time.Now()
	var proxy *url.URL
	var err error
	if task.Proxy != "" {
		proxy, err = url.Parse(task.Proxy)
		if err != nil {
			task.Error = err
			return task
		}
	}
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	if proxy != nil {
		client.Transport = &http.Transport{
			Proxy:             http.ProxyURL(proxy),
			DisableKeepAlives: true,
		}
	}
	req, err := http.NewRequest(http.MethodGet, task.Hostname, nil)
	if err != nil {
		task.Error = err
		return task
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:58.0) Gecko/20100101 Firefox/58.0")
	req.Header.Set("Connection", "close")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Referer", "https://www.google.com/")
	resp, err := client.Do(req)
	if err != nil {
		if resp != nil {
			_ = resp.Body.Close()
		}
		task.Error = err
		return task
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		task.Error = err
		err = resp.Body.Close()
		if err != nil {
			task.Error = err
		}
		return task
	}
	task.Body = body
	task.ResponseTime = time.Since(startTime)
	err = resp.Body.Close()
	if err != nil {
		task.Error = err
	}
	return task
}
