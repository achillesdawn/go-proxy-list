package proxyscrape

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"proxy-list/proxylist/proxy"
	"sync"
)

func getProxyScrapeData() (*proxyScrapeResponse, error) {
	targetUrl := "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&proxy_format=protocolipport&format=json"

	bytesData, err := proxy.Request(targetUrl)
	if err != nil {
		return nil, err
	}

	var data proxyScrapeResponse
	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		return nil, err
	}

	slog.Info("proxy scrape", slog.Int("total records", data.TotalRecords))

	return &data, nil
}

func WorkingProxies() ([]*ProxyScrapeProxy, error) {

	data, err := getProxyScrapeData()
	if err != nil {
		return nil, fmt.Errorf("could not get proxy data: %w", err)
	}

	var workingChan = make(chan *ProxyScrapeProxy, 100)
	var discarded = 0

	waitGroup := sync.WaitGroup{}

	for _, proxy := range data.Proxies {

		if !proxy.Ssl {
			discarded += 1
			continue
		}

		if proxy.Protocol == "socks5" || proxy.Protocol == "socks4" {

			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()

				ok, _ := proxy.TestProxy()

				if ok {
					workingChan <- &proxy
				}
			}()
		} else {
			slog.Warn(
				"[proxy scrape] protocol not supoorted",
				slog.String("protocol", proxy.Protocol),
			)
		}
	}

	waitGroup.Wait()

	close(workingChan)

	working := make([]*ProxyScrapeProxy, 0)

	for item := range workingChan {
		working = append(working, item)
	}

	return working, nil
}
