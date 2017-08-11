package pool

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

func crawl(t Task) Task {
	startTime := time.Now()
	client := &http.Client{
		Timeout: timeout,
	}
	if t.Proxy.Host != "" {
		client.Transport = &http.Transport{
			Proxy:             http.ProxyURL(t.Proxy),
			DisableKeepAlives: true,
		}
	}
	req, err := http.NewRequest("GET", t.Hostname, nil)
	if err != nil {
		t.Error = err
		return t
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		t.Error = err
		return t
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error = err
		return t
	}
	err = resp.Body.Close()
	if err != nil {
		t.Error = err
		return t
	}
	t.Body = body
	t.Response = resp
	t.ResponceTime = time.Since(startTime)
	return t
}
