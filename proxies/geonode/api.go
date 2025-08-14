package geonode

import (
	"time"
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

func (p *Proxy) ProxyIP() string {
	return p.IP
}
