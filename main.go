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
)

func TestUpwork(client *http.Client) error {

	url := "https://www.upwork.com/nx/search/jobs/?q=golang"

	headers := map[string]string{
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language":           "en-US,en;q=0.9",
		"cache-control":             "no-cache",
		"pragma":                    "no-cache",
		"priority":                  "u=0, i",
		"sec-ch-ua":                 "\"Not)A;Brand\";v=\"99\", \"Google Chrome\";v=\"127\", \"Chromium\";v=\"127\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Linux\"",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		cancel()
		panic(err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute request: %w", err)
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
