package proxifly

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetProxies(t *testing.T) {
	proxies, err := getProxies()
	require.NoError(t, err)

	for _, proxy := range proxies {
		indented, err := json.MarshalIndent(proxy, "", "\t")
		require.NoError(t, err)

		fmt.Println(string(indented))
	}
}

func TestWorkingProxies(t *testing.T) {
	proxies, err := WorkingProxies()
	require.NoError(t, err)

	for proxy := range proxies {
		fmt.Println(proxy.Protocol, proxy.Geolocation)

	}
}
