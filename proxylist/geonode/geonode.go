package geonode

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"proxy-list/proxylist/common"
	"sync"
)

func socksUrl(protocol, ip, port string) string {
	return fmt.Sprintf("%s://%s:%s", protocol, ip, port)
}

func (g *GeonodeProxy) CreateClient() (*http.Client, error) {

	for _, protocol := range g.Protocols {

		switch protocol {
		case common.ProtocolSocks4:
			return g.CreateSocks4Client(), nil
		case common.ProtocolSocks5:
			return g.CreateSocks5Client(), nil
		}
	}

	return nil, fmt.Errorf("no valid protocol found")
}

func pageURL(page uint8) string {
	return fmt.Sprintf(
		"https://proxylist.geonode.com/api/proxy-list?lpage=%d&limit=500&sort_by=lastChecked&sort_type=desc",
		page,
	)
}

func collectProxies() (map[string][]GeonodeProxy, error) {

	var currentPage uint8 = 1
	var count = 0

	// map of protocol to proxies that work for that protocol, so 3 is the expected capacity
	results := make(map[string][]GeonodeProxy, 3)

	for {
		page := pageURL(currentPage)

		byteData, err := common.Request(page)
		if err != nil {
			return nil, fmt.Errorf("geonodes request error: %w", err)
		}

		var data geonodeResponse

		err = json.Unmarshal(byteData, &data)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal: %w", err)
		}

		for _, proxy := range data.Data {

			for _, protocol := range proxy.Protocols {
				value, exists := results[protocol]

				if !exists {
					value = make([]GeonodeProxy, 0)
				}

				value = append(value, proxy)
				results[protocol] = value
			}
		}

		count += len(data.Data)

		slog.Info("geonode collecting", slog.Int("count", count))

		if count >= data.Total {
			break
		}
	}

	return results, nil
}

func WorkingProxies() ([]*GeonodeProxy, error) {
	proxies, err := collectProxies()
	if err != nil {
		return nil, fmt.Errorf("geonodes collect proxy: %w", err)
	}

	workingChan := make(chan *GeonodeProxy, 100)

	waitGroup := sync.WaitGroup{}

	for protocol, proxyList := range proxies {
		if protocol == common.ProtocolSocks4 || protocol == common.ProtocolSocks5 {
			for _, proxy := range proxyList {

				waitGroup.Add(1)

				go func() {
					defer waitGroup.Done()

					ok, _ := proxy.TestProxy()
					if ok {
						workingChan <- &proxy
					}
				}()
			}
		}
	}

	waitGroup.Wait()

	close(workingChan)

	working := make([]*GeonodeProxy, 0, 100)

	for proxy := range workingChan {
		working = append(working, proxy)
	}

	slog.Info(
		"geonodes working proxies",
		slog.Int("len", len(working)),
	)

	return working, nil
}
