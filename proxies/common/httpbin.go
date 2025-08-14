package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type httpbinIp struct {
	Origin string `json:"origin,omitempty"`
}

func TestProxy(client *http.Client, expectedIP string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, "https://httpbin.org/ip", nil)
	if err != nil {
		return false, fmt.Errorf("could not create test proxy request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("could not execute test proxy request: %w", err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("could not read test response: %w", err)
	}

	var data httpbinIp

	err = json.Unmarshal(b, &data)
	if err != nil {
		return false, fmt.Errorf("could not unmarshal httpbin response: %w", err)
	}

	if data.Origin == expectedIP {
		return true, nil
	} else {
		return false, nil
	}
}
