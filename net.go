package pool

import (
	"io/ioutil"
	"net/http"
	"time"
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
	req, err := http.NewRequest(http.MethodGet, t.Hostname, nil)
	if err != nil {
		t.Error = err
		return t
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
		t.Error = err
		return t
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error = err
		err = resp.Body.Close()
		if err != nil {
			t.Error = err
		}
		return t
	}
	t.Body = body
	t.ResponseTime = time.Since(startTime)
	err = resp.Body.Close()
	if err != nil {
		t.Error = err
	}
	return t
}
