package pool_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wwestgarth/liq-pool-amm/pool"
)

func almostEqual(u, v float64) bool {
	tol := 0.0000001
	fmt.Println("equality check", u, v)
	return math.Abs(v-u) < tol
}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestConstantProduct(t *testing.T) {
	nAssets := float64(100)
	p := pool.NewConstantProductPool(nAssets, nAssets, nAssets*nAssets)

	// what are the current prices
	a, b := p.Prices()
	require.True(t, almostEqual(a, 1))
	require.True(t, almostEqual(b, 1))

	// I want to sell 25 X, how many Y will I get
	dy, err := p.Trade(50, pool.SideSell)
	require.NoError(t, err)

	// dy = y*dx/(x + dx)
	// dy = 150 * 25 / 150 + 25 =
	require.NoError(t, p.Verify())
	require.True(t, almostEqual(dy, 33.33333333333333))
}
