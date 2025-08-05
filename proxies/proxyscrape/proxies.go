package proxyscrape

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/achillesdawn/proxy-list/proxies/common"
)

func proxyScrapeJSON() (*proxyScrapeResponse, error) {
	targetURL := "https://api.proxyscrape.com/v4/free-proxy-list/get?request=display_proxies&proxy_format=protocolipport&format=json"

	bytesData, err := common.Request(targetURL)
	if err != nil {
		return nil, err
	}

	var data proxyScrapeResponse

	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		return nil, err
	}

	slog.Info("proxy scrape",
		slog.Int("got records", data.ShownRecords),
		slog.Int("total records", data.TotalRecords),
	)

	return &data, nil
}

func WorkingProxies() ([]*Proxy, error) {

	data, err := proxyScrapeJSON()
	if err != nil {
		return nil, fmt.Errorf("could not get proxy data: %w", err)
	}

	var workingChan = make(chan *Proxy, 100)

	var discarded = 0

	waitGroup := sync.WaitGroup{}

	for _, proxy := range data.Proxies {

		if !proxy.Ssl {
			discarded += 1
			continue
		}

		if proxy.Protocol == common.ProtocolSocks5 || proxy.Protocol == common.ProtocolSocks4 {

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
				"[proxy scrape] protocol not supported",
				slog.String("protocol", proxy.Protocol),
			)
		}
	}

	waitGroup.Wait()

	close(workingChan)

	working := make([]*Proxy, 0)

	for item := range workingChan {
		working = append(working, item)
	}

	slog.Info(
		"working proxies",
		slog.Int("len", len(working)),
		slog.Int("discarded", discarded),
		slog.Int("total", len(data.Proxies)),
	)

	return working, nil
}
