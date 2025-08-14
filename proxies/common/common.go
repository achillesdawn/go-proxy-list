package common

import (
	"bufio"
	"context"
	"crypto/tls"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"h12.io/socks"
)

type Clientable interface {
	CreateClient() (*http.Client, error)
}

type protocol = string

const (
	ProtocolSocks4 protocol = "socks4"
	ProtocolSocks5 protocol = "socks5"
	ProtocolHTTP   protocol = "http"
	ProtocolHTTPS  protocol = "https"
)

func HttpProxyClient(proxyUrl *url.URL) *http.Client {
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

func Request(URL string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*60,
	)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
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
