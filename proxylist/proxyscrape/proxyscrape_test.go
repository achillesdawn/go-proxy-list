package proxyscrape

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorkingProxies(t *testing.T) {
	_, err := WorkingProxies()
	require.NoError(t, err)
}
