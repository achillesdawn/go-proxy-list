package geonode

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/achillesdawn/proxy-list/proxies/common"
)

func pageURL(page uint8, country string) string {
	if country == "" {
		return fmt.Sprintf(
			"https://proxylist.geonode.com/api/proxy-list?lpage=%d&limit=500&sort_by=lastChecked&sort_type=desc",
			page,
		)
	} else {
		return fmt.Sprintf(
			"https://proxylist.geonode.com/api/proxy-list?lpage=%d&country=%s&limit=500&sort_by=lastChecked&sort_type=desc",
			page,
			country,
		)
	}
}

func collectProxies(country string) (map[string][]Proxy, error) {

	var currentPage uint8 = 1
	var count = 0

	// map of protocol to proxies that work for that protocol, so 3 is the expected capacity
	results := make(map[string][]Proxy, 3)

	for {
		page := pageURL(currentPage, country)

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
					value = make([]Proxy, 0)
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

func checkProxies(proxies map[string][]Proxy) (<-chan Proxy, func(), error) {

	proxyList := make([]Proxy, 0)

	for _, value := range proxies {
		proxyList = append(proxyList, value...)
	}

	return common.TestProxies(proxyList)
}

func WorkingProxiesCountry(country string) (<-chan Proxy, func(), error) {
	proxies, err := collectProxies(country)
	if err != nil {
		return nil, nil, fmt.Errorf("geonodes collect proxy: %w", err)
	}

	return checkProxies(proxies)
}

func WorkingProxies() (<-chan Proxy, func(), error) {
	proxies, err := collectProxies("")
	if err != nil {
		return nil, nil, fmt.Errorf("geonodes collect proxy: %w", err)
	}
	return checkProxies(proxies)
}
