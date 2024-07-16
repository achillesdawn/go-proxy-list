package main

import (
	"fmt"
	"proxy-test/proxy"
	"sync"
)

func main() {
	proxies, err := proxy.CheckGeoNodes()
	if err != nil {
		panic(err)
	}

	waitGroup := sync.WaitGroup{}

	working := 0
	for protocol, proxyList := range proxies {

		if protocol == "socks5" || protocol == "socks4" {
			for _, proxy := range proxyList {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					ok, err := proxy.TestProxy()
					if err != nil {

					}
					if ok {
						working += 1
					}
				}()
			}
		}
	}

	waitGroup.Wait()
	fmt.Println("working:", working)
}
