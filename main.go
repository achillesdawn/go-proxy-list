package main

import (
	"log/slog"
	"proxy-list/proxylist/geonode"
	"proxy-list/proxylist/proxyscrape"
)

func main() {
	proxies, err := proxyscrape.WorkingProxies()
	if err != nil {

	}

	for _, p := range proxies {
		client := p.CreateClient()

		_ = client
	}

	geoproxies, err := geonode.WorkingProxies()
	if err != nil {
		panic(err)
	}

	for _, p := range geoproxies {

		client, err := p.CreateClient()
		if err != nil {
			slog.Error("creating client", "error", err)
			continue
		}

		_ = client

	}
}
