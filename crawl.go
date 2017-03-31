package pool

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type IpType struct {
	Addr        string
	Port        string
	Ssl         bool
	isWork      bool
	isAnon      bool
	ProxyChecks int
	CreateAt    time.Time
	LastCheck   time.Time
	Response    time.Duration
}

func crawl(t *task) {
	t.result, t.err = fetchBody(t.address)
}

func fetchBody(targetURL string) ([]byte, error) {
	client := &http.Client{
		// Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, err
}

func fetchBodyWProxy(targetURL string, proxy IpType) ([]byte, error) {
	client := &http.Client{
	// Timeout: time.Duration(timeout) * time.Second,
	}
	if proxy.Addr != "" {
		var (
			proxyURL *url.URL
			err      error
		)
		if proxy.Ssl {
			proxyURL, err = url.Parse("https://" + proxy.Addr + ":" + proxy.Port)
		} else {
			proxyURL, err = url.Parse("http://" + proxy.Addr + ":" + proxy.Port)
		}
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			Proxy:             http.ProxyURL(proxyURL),
			DisableKeepAlives: true,
		}
	}
	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, err
}
