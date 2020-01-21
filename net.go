package pool

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// func randomUA() string {
// 	userAgents := []string{
// 		"Mozilla/5.0 (Linux; Android 8.0.0; SM-G960F Build/R16NW) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.84 Mobile Safari/537.36",
// 		"Mozilla/5.0 (Linux; Android 7.0; SM-G892A Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/60.0.3112.107 Mobile Safari/537.36",
// 		"Mozilla/5.0 (Linux; Android 7.0; SM-G930VC Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/58.0.3029.83 Mobile Safari/537.36",
// 		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 6P Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.83 Mobile Safari/537.36",
// 		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1",
// 		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15A5341f Safari/604.1",
// 		"Mozilla/5.0 (Apple-iPhone7C2/1202.466; U; CPU like Mac OS X; en) AppleWebKit/420+ (KHTML, like Gecko) Version/3.0 Mobile/1A543 Safari/419.3",
// 		"Mozilla/5.0 (Windows Phone 10.0; Android 6.0.1; Microsoft; RM-1152) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Mobile Safari/537.36 Edge/15.15254",
// 		"Mozilla/5.0 (Linux; Android 7.0; SM-T827R4 Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.116 Safari/537.36",
// 		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
// 		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
// 		"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
// 		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
// 		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
// 		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:64.0) Gecko/20100101 Firefox/64.0",
// 		"Mozilla/5.0 (X11; Linux i686; rv:64.0) Gecko/20100101 Firefox/64.0",
// 		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1",
// 		"Links (2.7; Linux 3.7.9-2-ARCH x86_64; GNU C 4.7.1; text)",
// 		"Lynx/2.8.8dev.3 libwww-FM/2.14 SSL-MM/1.4.1",
// 		"Opera/9.80 (X11; Linux i686; Ubuntu/14.10) Presto/2.12.388 Version/12.16",
// 		"Opera/9.80 (Windows NT 6.0) Presto/2.12.388 Version/12.14",
// 		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
// 		"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
// 		"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)",
// 		"DuckDuckBot/1.0; (+http://duckduckgo.com/duckduckbot.html)",
// 		"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
// 		"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
// 		"Sogou Pic Spider/3.0( http://www.sogou.com/docs/help/webmasters.htm#07)",
// 		"Sogou head spider/3.0( http://www.sogou.com/docs/help/webmasters.htm#07)",
// 		"Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)",
// 		"Sogou Orion spider/3.0( http://www.sogou.com/docs/help/webmasters.htm#07)",
// 		"Sogou-Test-Spider/4.0 (compatible; MSIE 5.5; Windows 98)",
// 		"Mozilla/5.0 (compatible; Konqueror/3.5; Linux) KHTML/3.5.5 (like Gecko) (Exabot-Thumbnails)",
// 		"Mozilla/5.0 (compatible; Exabot/3.0; +http://www.exabot.com/go/robot)",
// 		"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
// 		"ia_archiver (+http://www.alexa.com/site/help/webmasters; crawler@alexa.com)",
// 	}
// 	return userAgents[rand.Intn(len(userAgents))]
// }

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
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
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
