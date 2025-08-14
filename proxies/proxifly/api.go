package proxifly

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/achillesdawn/proxy-list/proxies/common"
)

type Proxy struct {
	Proxy       string      `json:"proxy"`
	Protocol    string      `json:"protocol"`
	IP          string      `json:"ip"`
	Port        int         `json:"port"`
	HTTPS       bool        `json:"https"`
	Anonymity   string      `json:"anonymity"`
	Score       int         `json:"score"`
	Geolocation Geolocation `json:"geolocation"`
}

type Geolocation struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

func (p *Proxy) CreateClient() (*http.Client, error) {
	var client *http.Client
	var err error

	switch p.Protocol {
	case common.ProtocolSocks4:
		client = p.createSocks4Client()
	case common.ProtocolSocks5:
		client, err = p.createSocks5Client()
		if err != nil {
			return nil, err
		}
	case common.ProtocolHTTP, common.ProtocolHTTPS:
		client, err = p.createHTTPClient()
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("protocol not supported: %s", p.Protocol))
	}

	return client, nil
}

func (p *Proxy) createSocks5Client() (*http.Client, error) {
	URL, err := url.Parse(p.Proxy)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %s %w", p.Proxy, err)
	}
	return common.Socks5Client(URL), nil
}

func (p *Proxy) createSocks4Client() *http.Client {
	return common.Socks4Client(p.Proxy)
}

func (p *Proxy) createHTTPClient() (*http.Client, error) {
	URL, err := url.Parse(p.Proxy)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %s %w", p.Proxy, err)
	}
	return common.HttpProxyClient(URL), nil
}
