package proxyscrape

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/achillesdawn/proxy-list/proxies/common"
)

func generateURL(countries []string) string {
	if len(countries) > 0 {
		s := strings.Join(countries, ",")
		return fmt.Sprintf(
			"https://api.proxyscrape.com/v4/free-proxy-list/get?request=display_proxies&country=%s&proxy_format=protocolipport&format=json",
			s,
		)
	} else {
		return "https://api.proxyscrape.com/v4/free-proxy-list/get?request=display_proxies&proxy_format=protocolipport&format=json"
	}
}

func proxyScrapeJSON(countries []string) (*proxyScrapeResponse, error) {

	targetURL := generateURL(countries)

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

func WorkingProxies() (<-chan Proxy, func(), error) {

	data, err := proxyScrapeJSON([]string{})
	if err != nil {
		return nil, nil, fmt.Errorf("could not get proxy data: %w", err)
	}

	return common.TestProxies(data.Proxies)
}

func WorkingProxiesCountries(countries []string) (<-chan Proxy, func(), error) {

	data, err := proxyScrapeJSON(countries)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get proxy data: %w", err)
	}

	return common.TestProxies(data.Proxies)
}
