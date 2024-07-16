package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func writeToFile(res *http.Response) {
	file, err := os.Create("result.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	n, err := io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("wrote", n, "bytes")
}

func main() {
	targetUrl := "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&proxy_format=protocolipport&format=json"

	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(err.Error())
	}

	reader := io.LimitReader(bufio.NewReader(res.Body), 984735+100_000)
	b, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	var data ProxyResponse
	err = json.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}

	fmt.Println("got", data.TotalRecords)

	var counts = make(map[string][]*Proxy)

	var discared = 0

	for _, proxy := range data.Proxies {

		if !proxy.Ssl {
			discared += 1
			continue
		}

		if proxy.Protocol == "socks5" {
			err := proxy.TestProxy()
			if err != nil {
				slog.Error(err.Error())
				continue
			}

		}

		value, exists := counts[proxy.Protocol]
		if !exists {
			counts[proxy.Protocol] = make([]*Proxy, 0)
			value = counts[proxy.Protocol]
		}
		value = append(value, &proxy)
		counts[proxy.Protocol] = value
	}

	for key, value := range counts {
		fmt.Println(key, "count:", len(value))
	}

	fmt.Println("discarded", discared)
}
