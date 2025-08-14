package proxifly

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
