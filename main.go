package main

import (
	"fmt"
	"proxy-test/proxy"
)

func main() {
	proxies, err := proxy.CheckGeoNodes()
	if err != nil {
		panic(err)
	}

	for protocol, proxy := range proxies {
		fmt.Println(protocol, len(proxy))
	}

}
