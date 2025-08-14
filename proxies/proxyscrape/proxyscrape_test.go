package proxyscrape

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkingProxies(t *testing.T) {
	proxies, cancel, err := WorkingProxies()
	require.NoError(t, err)

	for proxy := range proxies {

		fmt.Println("working", proxy.Proxy)
		cancel()
	}
}
