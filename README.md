# proxy-list
Simplified API to get and use free proxies from ProxyScrape and GeoNodes

## Getting started

```bash
go get github.com/achillesdawn/proxy-list
```

## Usage

```go
package main

import (
	"log/slog"

	"github.com/achillesdawn/proxy-list/proxies/geonode"
	"github.com/achillesdawn/proxy-list/proxies/proxyscrape"
)

func main() {
	proxies, err := proxyscrape.WorkingProxies()
	if err != nil {
		// handle error
		// for now, we shall panic
		panic(err)
	}

	for _, proxy := range proxies {
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

	for _, proxy := range geoproxies {
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

```
