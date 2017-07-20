package pool

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

func crawl(t *Task) {
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
	req, err := http.NewRequest("GET", t.Target.String(), nil)
	if err != nil {
		t.Error = err
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		t.Error = err
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error = err
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Error = err
		return
	}
	t.Body = body
	t.Response = resp
	t.ResponceTime = time.Since(startTime)
}
