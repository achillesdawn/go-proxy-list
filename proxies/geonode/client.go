package geonode

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/achillesdawn/proxy-list/proxies/common"
)

func socksUrl(protocol, ip, port string) string {
	return fmt.Sprintf("%s://%s:%s", protocol, ip, port)
}

func (g Proxy) CreateClient() (*http.Client, error) {

	for _, protocol := range g.Protocols {

		switch protocol {
		case common.ProtocolSocks4:
			return g.createSocks4Client(), nil
		case common.ProtocolSocks5:
			return g.createSocks5Client(), nil
		case common.ProtocolHTTP, common.ProtocolHTTPS:
			URL := socksUrl(protocol, g.IP, g.Port)
			return g.createHTTPClient(URL)
		}
	}

	return nil, fmt.Errorf("no valid protocol found")
}

func (g Proxy) createSocks5Client() *http.Client {
	proxyUrl := socksUrl(common.ProtocolSocks5, g.IP, g.Port)

	URL, err := url.Parse(proxyUrl)
	if err != nil {
		panic(err)
	}

	return common.Socks5Client(URL)
}

func (g Proxy) createSocks4Client() *http.Client {
	proxyUrl := socksUrl(common.ProtocolSocks4, g.IP, g.Port)

	return common.Socks4Client(proxyUrl)
}

func (p Proxy) createHTTPClient(URL string) (*http.Client, error) {
	parsed, err := url.Parse(URL)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %s %w", URL, err)
	}
	return common.HttpProxyClient(parsed), nil
}
