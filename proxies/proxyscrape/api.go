package proxyscrape

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"proxy-list/proxies/common"
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

func (p *ProxyScrapeProxy) CreateClient() *http.Client {
	var client *http.Client

	switch p.Protocol {
	case common.ProtocolSocks4:
		client = p.CreateSocks4Client()
	case common.ProtocolSocks5:
		client = p.CreateSocks5Client()
	default:
		panic(fmt.Sprintf("protocol not supported: %s", p.Protocol))
	}

	return client
}

func (p *ProxyScrapeProxy) CreateSocks5Client() *http.Client {
	URL, err := url.Parse(p.Proxy)
	if err != nil {
		panic(err)
	}
	return common.Socks5Client(URL)
}

func (p *ProxyScrapeProxy) CreateSocks4Client() *http.Client {
	return common.Socks4Client(p.Proxy)
}

func (p *ProxyScrapeProxy) TestProxy() (bool, error) {

	slog.Info(
		"[proxy scrape] testing URL",
		slog.String("URL", p.Proxy),
	)

	var client *http.Client

	switch p.Protocol {
	case common.ProtocolSocks4:
		client = p.CreateSocks4Client()
	case common.ProtocolSocks5:
		client = p.CreateSocks5Client()
	default:
		return false, fmt.Errorf("expecting protocol either socks4 or socks5")
	}

	ok, _ := common.TestProxy(client, p.Ip)

	if ok {
		return true, nil
	} else {
		return false, nil
	}
}
