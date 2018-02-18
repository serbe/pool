package pool

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/headzoo/surf"
)

func (p *Pool) crawl(t *Task) *Task {
	startTime := time.Now()
	client := &http.Client{
		Timeout: p.timeout,
	}
	if t.Proxy != nil {
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
	// ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	// defer cancel()
	// req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		t.Error = err
		return t
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error = err
		_ = resp.Body.Close()
		return t
	}
	t.Body = body
	t.Response = resp
	t.ResponceTime = time.Since(startTime)
	_ = resp.Body.Close()
	return t
}

func (p *Pool) surfCrawl(t *Task) *Task {
	startTime := time.Now()
	bow := surf.NewBrowser()
	err := bow.Open(t.Hostname)
	if err != nil {
		t.Error = err
		return t
	}
	// t.Response = bow.ResponseHeaders()
	t.ResponceTime = time.Since(startTime)
	t.Browser = bow
	return t
}
