package pool

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func crawl(t *task) {
	t.result, t.err = fetchBody(t.address, t.proxy)
}

func fetchBody(targetURL string, proxy string) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy:             http.ProxyURL(proxyURL),
				DisableKeepAlives: true,
			}
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
