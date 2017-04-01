package pool

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func crawl(t *Task) error {
	startTime := time.Now()
	client := &http.Client{
		Timeout: timeout,
	}
	if t.Proxy != "" {
		proxyURL, err := url.Parse(t.Proxy)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy:             http.ProxyURL(proxyURL),
				DisableKeepAlives: true,
			}
		}
	}
	resp, err := client.Get(t.Address)
	if err != nil {
		t.Error = err
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error = err
		return err
	}
	t.Body = body
	t.ResponceTime = time.Now().Sub(startTime)
	return err
}
