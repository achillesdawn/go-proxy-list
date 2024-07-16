package main

import (
	"fmt"
	"proxy-test/proxy"
)

func main() {
	proxies, err := proxy.CheckProxyScrape()
	if err != nil {
		panic(err)
	}

	fmt.Println("working:", len(proxies))
}
