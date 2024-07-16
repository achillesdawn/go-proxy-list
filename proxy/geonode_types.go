package proxy

import (
	"encoding/json"
	"fmt"
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

	GeonodeResponse struct {
		Data  []GeonodeProxy `json:"data"`
		Total int            `json:"total"`
		Page  int            `json:"page"`
		Limit int            `json:"limit"`
	}
)

func urlForPage(page uint8) string {
	return fmt.Sprintf("https://proxylist.geonode.com/api/proxy-list?lpage=%d&limit=500&sort_by=lastChecked&sort_type=desc", page)
}

func CheckGeoNodes() (map[string][]GeonodeProxy, error) {

	var currentPage uint8 = 1
	var count int = 0

	results := make(map[string][]GeonodeProxy, 3)

	for {
		page := urlForPage(currentPage)
		byteData, err := makeRequest(page)
		if err != nil {
			panic(err)
		}

		var data GeonodeResponse
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
