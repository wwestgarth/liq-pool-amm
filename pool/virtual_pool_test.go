package pool_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wwestgarth/liq-pool-amm/pool"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestVirtualPool(t *testing.T) {
	upper := 15.0
	base := 10.0
	lower := 5.0

	commitment := 100000.0

	p := pool.NewConcentratedLiquidityPool(base, lower, upper, commitment)
	require.InDelta(t, 10.0, p.FairPrice(), 0.0001)

	// for a small position should it be close to the base price??? what am I missing
	p.Trade(pool.SideSell, 1)
	require.InDelta(t, 10.0, p.FairPrice(), 0.0001)

	// short position of 1, fair price is now 27.5, diff from base = 17.5
	// long position of 1,  fair price is now 4.3,  diff from base = 5.7

	fmt.Println("VolBetween", p.VolumeBetweenPrices(5, 10))
	fmt.Println("VolBetween", p.VolumeBetweenPrices(10, 15))
	require.False(t, true)

}
