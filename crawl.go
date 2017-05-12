package pool

import (
	"io/ioutil"
	"net/http"
	"time"
)

func crawl(t *Task) error {
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
	resp, err := client.Get(t.Target.String())
	if err != nil {
		t.Error = err
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error = err
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		t.Error = err
		return err
	}
	t.Body = body
	t.Response = resp
	t.ResponceTime = time.Since(startTime)
	return err
}
