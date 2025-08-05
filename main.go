package main

import (
	"log/slog"

	"github.com/achillesdawn/proxy-list/proxies/geonode"
	"github.com/achillesdawn/proxy-list/proxies/proxyscrape"
)

func main() {
	// proxy is a <-chan of proxies
	proxies, err := proxyscrape.WorkingProxies()
	if err != nil {
		// handle error
		// for now, we shall panic
		panic(err)
	}

	// we can iterate over the channel to get proxies as they come after being tested
	for proxy := range proxies {
		// client is of type *http.Client, already configured with proxy information
		// it can be used simply to execute requests
		// client.Do(request)
		client := proxy.CreateClient()

		_ = client
	}

	geoproxies, err := geonode.WorkingProxies()
	if err != nil {
		panic(err)
	}

	for proxy := range geoproxies {
		// client is of type *http.Client, already configured with proxy information
		// it can be used simply to execute requests
		// client.Do(request)
		client, err := proxy.CreateClient()
		if err != nil {
			slog.Error("creating client", "error", err)
			continue
		}

		_ = client

	}
}
