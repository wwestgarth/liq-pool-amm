package pool_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wwestgarth/liq-pool-amm/pool"
)

func scaleUp(s uint64) uint64 {
	return s * 10000000
}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestConstantProduct(t *testing.T) {
	nAssets := scaleUp(uint64(100))
	p := pool.NewConstantProductPool(nAssets, nAssets, nAssets*nAssets)

	// I want to sell 25 X, how many Y will I get
	dy, err := p.Trade(scaleUp(50), pool.SideSell)
	require.NoError(t, err)

	// dy = y*dx/(x + dx)
	// dy = 150 * 25 / 150 + 25 =
	require.NoError(t, p.Verify())
	require.Equal(t, dy, 21.428571428571427)
}
