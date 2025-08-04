package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"proxy-list/proxylist/geonode"
	"proxy-list/proxylist/proxyscrape"
	"time"

	"golang.org/x/net/http2"
)

func withAnonHeaders(req *http.Request) {
	anonHeaders := map[string]string{
		"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64; rv:135.0) Gecko/20100101 Firefox/135.0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.5",
		"Sec-GPC":                   "1",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "cross-site",
		"Priority":                  "u=0, i",
		"Pragma":                    "no-cache",
		"Cache-Control":             "no-cache",
	}

	for key, value := range anonHeaders {
		req.Header.Add(key, value)
	}
}

func TestUpwork(client *http.Client) error {

	url := "https://www.upwork.com/nx/search/jobs/?q=golang"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		cancel()
		panic(err)
	}

	client.Transport = &http2.Transport{}

	withAnonHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute request: %w", err)
	}

	if res.StatusCode > 300 {
		return fmt.Errorf("status: %s", res.Status)
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read body: %w", err)
	}

	fmt.Println(string(bytes))
	return nil
}

func main() {
	proxies, err := proxyscrape.WorkingProxies()
	if err != nil {
		panic(err)
	}

	for _, p := range proxies {
		client := p.CreateClient()

		err := TestUpwork(client)
		if err != nil {
			slog.Error("fail", "error", err)
		}
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

		err = TestUpwork(client)
		if err != nil {
			slog.Error("fail", "error", err)
		}
	}
}
