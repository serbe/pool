package pool

import (
	"io/ioutil"
	"net/http"
)

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

// func fetchBodyWProxy(targetURL string, proxy ipType) ([]byte, error) {
// 	client := &http.Client{
// 		Timeout: time.Duration(timeout) * time.Second,
// 	}
// 	if proxy.Addr != "" {
// 		var (
// 			proxyURL *url.URL
// 			err      error
// 		)
// 		if proxy.Ssl {
// 			proxyURL, err = url.Parse("https://" + proxy.Addr + ":" + proxy.Port)
// 		} else {
// 			proxyURL, err = url.Parse("http://" + proxy.Addr + ":" + proxy.Port)
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		client.Transport = &http.Transport{
// 			Proxy:             http.ProxyURL(proxyURL),
// 			DisableKeepAlives: true,
// 		}
// 	}
// 	resp, err := client.Get(targetURL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	return body, err
// }
