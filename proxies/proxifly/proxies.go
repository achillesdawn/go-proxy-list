package proxifly

import (
	"encoding/json"
	"fmt"

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

func WorkingProxies() (<-chan Proxy, func(), error) {
	proxies, err := getProxies()
	if err != nil {
		return nil, nil, err
	}
	return common.TestProxies(proxies)
}
