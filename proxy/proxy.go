package proxy

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"h12.io/socks"
)

type HttpbinIp struct {
	Origin string `json:"origin,omitempty"`
}

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

func createSocks5Client(proxyUrl *url.URL) *http.Client {
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

func createSocks4Client(proxyUrl string) *http.Client {

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

func testProxy(client *http.Client, ip string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, "https://httpbin.org/ip", nil)
	if err != nil {
		return false, err
	}
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	fmt.Println(string(b))

	var data HttpbinIp
	err = json.Unmarshal(b, &data)
	if err != nil {
		return false, err
	}

	if data.Origin == ip {
		return true, nil
	} else {
		return false, nil
	}
}

func makeRequest(targetUrl string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	reader := io.LimitReader(bufio.NewReader(res.Body), 984735+100_000)
	b, err := io.ReadAll(reader)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return b, nil
}
