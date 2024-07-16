package proxy

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
)

type (
	IpData struct {
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

	Proxy struct {
		Alive          bool    `json:"alive,omitempty"`
		AliveSince     float64 `json:"alive_since,omitempty"`
		Anonimity      string  `json:"anonimity,omitempty"`
		AverageTimeout float32 `json:"average_timeout,omitempty"`
		FirstSeen      float64 `json:"first_seen,omitempty"`
		IpData         IpData  `json:"ip_data,omitempty"`
		LastSeen       float64 `json:"last_seen,omitempty"`
		Port           int     `json:"port,omitempty"`
		Protocol       string  `json:"protocol,omitempty"`
		Proxy          string  `json:"proxy,omitempty"`
		Ssl            bool    `json:"ssl,omitempty"`
		Timeout        float64 `json:"timeout,omitempty"`
		TimesAlive     int     `json:"times_alive,omitempty"`
		TimesDead      int     `json:"times_dead,omitempty"`
		Uptime         float64 `json:"uptime,omitempty"`
		Ip             string  `json:"ip,omitempty"`
	}

	ProxyResponse struct {
		ShownRecords int     `json:"shown_records,omitempty"`
		TotalRecords int     `json:"total_records,omitempty"`
		Limit        int     `json:"limit,omitempty"`
		Skip         int     `json:"skip,omitempty"`
		Nextpage     bool    `json:"nextpage,omitempty"`
		Proxies      []Proxy `json:"proxies,omitempty"`
	}
)

func (p *Proxy) TestProxy() (bool, error) {

	slog.Info("testing", slog.String("url", p.Proxy))

	proxyUrl, err := url.Parse(p.Proxy)
	if err != nil {
		panic(err)
	}

	client := createSocks5Client(proxyUrl)

	ok, err := testProxy(client, p.Ip)

	if ok {
		return true, nil
	} else {
		return false, nil
	}
}

func ProxyScrape() (map[string][]*Proxy, error) {
	targetUrl := "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&proxy_format=protocolipport&format=json"

	bytesData, err := makeRequest(targetUrl)
	if err != nil {
		return nil, err
	}

	var data ProxyResponse
	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		return nil, err
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
			ok, err := proxy.TestProxy()
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if ok {

			}

		}

		value, exists := counts[proxy.Protocol]
		if !exists {
			value = make([]*Proxy, 0)
		}
		value = append(value, &proxy)
		counts[proxy.Protocol] = value
	}

	for key, value := range counts {
		fmt.Println(key, "count:", len(value))
	}

	fmt.Println("discarded", discared)
	return counts, nil
}
