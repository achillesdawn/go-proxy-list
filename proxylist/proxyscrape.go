package proxylist

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
)

type (
	proxyScrapeIpData struct {
		As            string  `json:"as,omitempty"`
		Asname        string  `json:"asname,omitempty"`
		City          string  `json:"city,omitempty"`
		Continent     string  `json:"continent,omitempty"`
		ContinentCode string  `json:"continentCode,omitempty"`
		Country       string  `json:"country,omitempty"`
		CountryCode   string  `json:"countryCode,omitempty"`
		District      string  `json:"district,omitempty"`
		Hosting       bool    `json:"hosting,omitempty"`
		Isp           string  `json:"isp,omitempty"`
		Lat           float32 `json:"lat,omitempty"`
		Lon           float32 `json:"lon,omitempty"`
		Mobile        bool    `json:"mobile,omitempty"`
		Org           string  `json:"org,omitempty"`
		Proxy         bool    `json:"proxy,omitempty"`
		RegionName    string  `json:"regionName,omitempty"`
		Status        string  `json:"status,omitempty"`
		Timezone      string  `json:"timezone,omitempty"`
		Zip           string  `json:"zip,omitempty"`
	}

	ProxyScrapeProxy struct {
		Alive          bool              `json:"alive,omitempty"`
		AliveSince     float64           `json:"alive_since,omitempty"`
		Anonimity      string            `json:"anonimity,omitempty"`
		AverageTimeout float32           `json:"average_timeout,omitempty"`
		FirstSeen      float64           `json:"first_seen,omitempty"`
		IpData         proxyScrapeIpData `json:"ip_data,omitempty"`
		LastSeen       float64           `json:"last_seen,omitempty"`
		Port           int               `json:"port,omitempty"`
		Protocol       string            `json:"protocol,omitempty"`
		Proxy          string            `json:"proxy,omitempty"`
		Ssl            bool              `json:"ssl,omitempty"`
		Timeout        float64           `json:"timeout,omitempty"`
		TimesAlive     int               `json:"times_alive,omitempty"`
		TimesDead      int               `json:"times_dead,omitempty"`
		Uptime         float64           `json:"uptime,omitempty"`
		Ip             string            `json:"ip,omitempty"`
	}

	proxyScrapeResponse struct {
		ShownRecords int                `json:"shown_records,omitempty"`
		TotalRecords int                `json:"total_records,omitempty"`
		Limit        int                `json:"limit,omitempty"`
		Skip         int                `json:"skip,omitempty"`
		Nextpage     bool               `json:"nextpage,omitempty"`
		Proxies      []ProxyScrapeProxy `json:"proxies,omitempty"`
	}
)

func (p *ProxyScrapeProxy) CreateSocks5Client() *http.Client {
	Url, err := url.Parse(p.Proxy)
	if err != nil {
		panic(err)
	}
	return createSocks5Client(Url)
}

func (p *ProxyScrapeProxy) CreateSocks4Client() *http.Client {
	return createSocks4Client(p.Proxy)
}

func (p *ProxyScrapeProxy) TestProxy() (bool, error) {

	slog.Info("testing", slog.String("url", p.Proxy))

	var client *http.Client

	if p.Protocol == "socks4" {
		client = p.CreateSocks4Client()
	} else if p.Protocol == "socks5" {
		client = p.CreateSocks5Client()
	}

	ok, _ := testProxy(client, p.Ip)

	if ok {
		return true, nil
	} else {
		return false, nil
	}
}

func getProxyScrapeData() (*proxyScrapeResponse, error) {
	targetUrl := "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&proxy_format=protocolipport&format=json"

	bytesData, err := makeRequest(targetUrl)
	if err != nil {
		return nil, err
	}

	var data proxyScrapeResponse
	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		return nil, err
	}

	fmt.Println("got", data.TotalRecords)
	return &data, nil
}

func ProxyScrapeWorkingProxies() ([]*ProxyScrapeProxy, error) {

	data, err := getProxyScrapeData()
	if err != nil {
		panic(err)
	}

	var workingChan = make(chan *ProxyScrapeProxy, 100)
	var discarded = 0

	waitGroup := sync.WaitGroup{}

	for _, proxy := range data.Proxies {

		if !proxy.Ssl {
			discarded += 1
			continue
		}

		if proxy.Protocol == "socks5" || proxy.Protocol == "socks4" {

			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()

				ok, _ := proxy.TestProxy()

				if ok {
					workingChan <- &proxy
				}
			}()
		}
	}

	waitGroup.Wait()
	close(workingChan)

	working := make([]*ProxyScrapeProxy, 0)

	for item := range workingChan {
		working = append(working, item)
	}

	return working, nil
}
