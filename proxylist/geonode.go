package proxylist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type (
	GeonodeProxy struct {
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
		Data  []GeonodeProxy `json:"data"`
		Total int            `json:"total"`
		Page  int            `json:"page"`
		Limit int            `json:"limit"`
	}
)

func socksUrl(protocol, ip, port string) string {
	return fmt.Sprintf("%s://%s:%s", protocol, ip, port)
}

func (g *GeonodeProxy) CreateSocks5Client() *http.Client {
	proxyUrl := socksUrl("socks5", g.IP, g.Port)
	Url, err := url.Parse(proxyUrl)
	if err != nil {
		panic(err)
	}
	return createSocks5Client(Url)
}

func (g *GeonodeProxy) CreateSocks4Client() *http.Client {
	proxyUrl := socksUrl("socks4", g.IP, g.Port)

	return createSocks4Client(proxyUrl)
}

func (g *GeonodeProxy) TestProxy() (bool, error) {

	// geonodes proxy has an array of protocols
	for _, protocol := range g.Protocols {
		if protocol == "socks5" {

			client := g.CreateSocks5Client()

			ok, err := testProxy(client, g.IP)
			if err != nil {
				return false, err
			}

			if ok {
				return true, nil
			}

		} else if protocol == "socks4" {

			client := g.CreateSocks4Client()

			ok, err := testProxy(client, g.IP)
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

func urlForPage(page uint8) string {
	return fmt.Sprintf("https://proxylist.geonode.com/api/proxy-list?lpage=%d&limit=500&sort_by=lastChecked&sort_type=desc", page)
}

func checkGeoNodes() (map[string][]GeonodeProxy, error) {

	var currentPage uint8 = 1
	var count int = 0

	results := make(map[string][]GeonodeProxy, 3)

	for {
		page := urlForPage(currentPage)
		byteData, err := makeRequest(page)
		if err != nil {
			panic(err)
		}

		var data geonodeResponse
		err = json.Unmarshal(byteData, &data)
		if err != nil {
			return nil, err
		}

		for _, proxy := range data.Data {

			for _, protocol := range proxy.Protocols {
				value, exists := results[protocol]

				if !exists {
					value = make([]GeonodeProxy, 0)
				}

				value = append(value, proxy)
				results[protocol] = value
			}
		}

		count += len(data.Data)

		fmt.Println("count", count)
		if count >= data.Total {
			break
		}
	}

	return results, nil
}

func GeonodesWorkingProxies() []*GeonodeProxy {
	proxies, err := checkGeoNodes()
	if err != nil {
		panic(err)
	}

	waitGroup := sync.WaitGroup{}

	workingChan := make(chan *GeonodeProxy, 100)

	for protocol, proxyList := range proxies {

		if protocol == "socks5" || protocol == "socks4" {
			for _, proxy := range proxyList {
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
	}

	waitGroup.Wait()
	close(workingChan)

	working := make([]*GeonodeProxy, 0, 100)
	for proxy := range workingChan {
		working = append(working, proxy)
	}
	fmt.Println("working:", len(working))
	return working
}
