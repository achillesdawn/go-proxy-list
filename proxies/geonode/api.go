package geonode

import (
	"net/http"
	"net/url"
	"time"

	"github.com/achillesdawn/proxy-list/proxies/common"
)

type (
	Proxy struct {
		ID                 string    `json:"_id"`
		IP                 string    `json:"ip"`
		AnonymityLevel     string    `json:"anonymityLevel"`
		Asn                string    `json:"asn"`
		City               string    `json:"city"`
		Country            string    `json:"country"`
		CreatedAt          time.Time `json:"created_at"`
		Google             bool      `json:"google"`
		Isp                string    `json:"isp"`
		LastChecked        int       `json:"lastChecked"`
		Latency            float64   `json:"latency"`
		Org                string    `json:"org"`
		Port               string    `json:"port"`
		Protocols          []string  `json:"protocols"`
		Speed              int       `json:"speed"`
		UpTime             float64   `json:"upTime"`
		UpTimeSuccessCount int       `json:"upTimeSuccessCount"`
		UpTimeTryCount     int       `json:"upTimeTryCount"`
		UpdatedAt          time.Time `json:"updated_at"`
		ResponseTime       int       `json:"responseTime"`
		Region             any       `json:"region,omitempty"`
		WorkingPercent     any       `json:"workingPercent,omitempty"`
	}

	geonodeResponse struct {
		Data  []Proxy `json:"data"`
		Total int     `json:"total"`
		Page  int     `json:"page"`
		Limit int     `json:"limit"`
	}
)

func (g *Proxy) createSocks5Client() *http.Client {
	proxyUrl := socksUrl(common.ProtocolSocks5, g.IP, g.Port)

	URL, err := url.Parse(proxyUrl)
	if err != nil {
		panic(err)
	}

	return common.Socks5Client(URL)
}

func (g *Proxy) createSocks4Client() *http.Client {
	proxyUrl := socksUrl(common.ProtocolSocks4, g.IP, g.Port)

	return common.Socks4Client(proxyUrl)
}

func (g *Proxy) TestProxy() (bool, error) {

	// geonodes proxy has an array of protocols
	for _, protocol := range g.Protocols {
		switch protocol {
		case common.ProtocolSocks5:
			c := g.createSocks5Client()

			ok, err := common.TestProxy(c, g.IP)
			if err != nil {
				return false, err
			}

			if ok {
				return true, nil
			}

		case common.ProtocolSocks4:
			c := g.createSocks4Client()

			ok, err := common.TestProxy(c, g.IP)
			if err != nil {
				return false, err
			}

			if ok {
				return true, nil
			}
		}
	}

	return false, nil
}
