package common

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"h12.io/socks"
)

type Clientable interface {
	CreateClient() (*http.Client, error)
	GetIP() string
	GetProtocol() string
	GetAddress() string
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

func TestProxies[T Clientable](proxies []T) (<-chan T, func(), error) {

	var workingChan = make(chan T, 100)
	var done = make(chan struct{})

	parentCtx, cancelAll := context.WithCancel(context.Background())

	go func() {
		defer cancelAll()

		<-done
	}()

	go func() {

		var waitGroup sync.WaitGroup

		for _, proxy := range proxies {

			client, err := proxy.CreateClient()
			if err != nil {
				panic(err)
			}

			waitGroup.Add(1)

			go func() {
				defer waitGroup.Done()

				ctx, cancel := context.WithTimeout(
					parentCtx,
					time.Second*60,
				)

				defer cancel()

				ok, err := TestProxy(ctx, client, proxy.GetIP())
				if errors.Is(err, context.Canceled) {
					return
				} else if err != nil {
					slog.Error("test fail",
						slog.String("protocol", proxy.GetProtocol()),
						slog.String("address", proxy.GetAddress()),
						slog.String("error", err.Error()))
				}

				if ok {
					workingChan <- proxy
				}
			}()
		}

		waitGroup.Wait()

		close(workingChan)
	}()

	return workingChan, func() { done <- struct{}{} }, nil
}
