```go
package main

import (
	"fmt"
	"proxy-list/proxylist"
)

func main() {
	proxyscrapeProxies, err := proxylist.ProxyScrapeWorkingProxies()
	if err != nil {
		panic(err)
	}

	fmt.Println("proxyscrape working:", len(proxyscrapeProxies))

	geonodeProxies := proxylist.GeonodesWorkingProxies()

	fmt.Println("geonodes working:", len(geonodeProxies))

}

```
