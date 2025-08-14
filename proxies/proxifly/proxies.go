package proxifly

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/achillesdawn/proxy-list/proxies/common"
)

func getProxies() ([]Proxy, error) {
	URL := "https://cdn.jsdelivr.net/gh/proxifly/free-proxy-list@main/proxies/all/data.json"

	b, err := common.Request(URL)
	if err != nil {
		return nil, fmt.Errorf("proxifly url get error: %w", err)
	}

	var proxies []Proxy

	err = json.Unmarshal(b, &proxies)
	if err != nil {
		return nil, fmt.Errorf("proxifly could not unmarshal json response: %w", err)
	}
	return proxies, nil
}

func WorkingProxies() (<-chan *Proxy, error) {
	proxies, err := getProxies()
	if err != nil {
		return nil, err
	}

	var workingChan = make(chan *Proxy, 100)

	go func() {

		var waitGroup sync.WaitGroup

		for _, proxy := range proxies {

			client, err := proxy.CreateClient()
			if err != nil {
				panic(err)
			}

			waitGroup.Add(1)

			go func() {
				defer waitGroup.Done()

				ok, err := common.TestProxy(client, proxy.IP)
				if err != nil {
					slog.Error("test fail",
						slog.String("protocol", proxy.Protocol),
						slog.String("address", proxy.Proxy),
						slog.String("error", err.Error()))
				}

				if ok {
					workingChan <- &proxy
				}
			}()
		}

		waitGroup.Wait()

		close(workingChan)
	}()

	return workingChan, nil
}
