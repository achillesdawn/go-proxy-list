package common

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"h12.io/socks"
)

const ProtocolSocks4 string = "socks4"
const ProtocolSocks5 string = "socks5"

func Socks5Client(proxyUrl *url.URL) *http.Client {
	client := &http.Client{
		Timeout: time.Millisecond * 10_000,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),

			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return client
}

func Socks4Client(proxyUrl string) *http.Client {

	dial := socks.Dial(proxyUrl)

	client := &http.Client{
		Timeout: time.Millisecond * 10_000,
		Transport: &http.Transport{
			Dial: dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return client
}

type httpbinIp struct {
	Origin string `json:"origin,omitempty"`
}

func Request(URL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info(
		"request status",
		slog.String("URL", URL),
		slog.String("status", res.Status),
	)

	reader := io.LimitReader(bufio.NewReader(res.Body), 1024*1024*5)

	b, err := io.ReadAll(reader)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return b, nil
}

func TestProxy(client *http.Client, ip string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, "https://httpbin.org/ip", nil)
	if err != nil {
		return false, fmt.Errorf("could not create test proxy request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("could not execute test proxy request: %w", err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("could not read test response: %w", err)
	}

	var data httpbinIp

	err = json.Unmarshal(b, &data)
	if err != nil {
		return false, fmt.Errorf("could not unmarshal httpbin response: %w", err)
	}

	if data.Origin == ip {
		return true, nil
	} else {
		return false, nil
	}
}
