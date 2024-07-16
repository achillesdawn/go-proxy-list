package main

import (
	"fmt"
	"proxy-test/proxylist"
)

func main() {
	proxies, err := proxylist.ProxyScrapeWorkingProxies()
	if err != nil {
		panic(err)
	}

	fmt.Println("working:", len(proxies))

}
