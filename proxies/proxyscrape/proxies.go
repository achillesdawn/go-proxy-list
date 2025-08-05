package proxyscrape

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"

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

func WorkingProxies() (<-chan *Proxy, error) {

	data, err := proxyScrapeJSON([]string{})
	if err != nil {
		return nil, fmt.Errorf("could not get proxy data: %w", err)
	}

	var workingChan = make(chan *Proxy, 100)

	go func() {

		waitGroup := sync.WaitGroup{}

		for _, proxy := range data.Proxies {

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
	}()

	return workingChan, nil
}

func WorkingProxiesCountries(countries []string) (<-chan *Proxy, error) {

	data, err := proxyScrapeJSON(countries)
	if err != nil {
		return nil, fmt.Errorf("could not get proxy data: %w", err)
	}

	var workingChan = make(chan *Proxy, 100)

	go func() {

		waitGroup := sync.WaitGroup{}

		for _, proxy := range data.Proxies {

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
	}()

	return workingChan, nil
}
