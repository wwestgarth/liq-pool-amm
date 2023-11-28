package pool

import "math"

type ConcentratedLiquidityPool struct {
	initial  float64 // initial commitment to the pool
	base     float64 // price where the pool has position 0, and contents it all physical tokens
	lower    float64 // price in [lower, base) pool will have a long position
	upper    float64 // price in (base, upper] pool will have a short position
	position float64 // the pools current position

	rf     float64 // risk factor
	lowerL float64
	upperL float64
}

func NewConcentratedLiquidityPool(b, l, u, i float64) *ConcentratedLiquidityPool {

	// just asssume risk factor is 1 for now while we get going, and so v_worst = initial commitment
	rf := 1.0
	vw := rf * i

	// calculate liqudity in each half
	// L = v_worst / ( sqrt(pu) - sqrt(pl) )

	// upper so interval [upper, base)
	upL := vw / (math.Sqrt(u) - math.Sqrt(b))

	// lower so interval (base, lower]
	lowL := vw / (math.Sqrt(b) - math.Sqrt(l))

	return &ConcentratedLiquidityPool{
		base:     b,
		lower:    l,
		upper:    u,
		initial:  i,
		position: 0.0,

		rf:     rf,
		lowerL: lowL,
		upperL: upL,
	}
}

// calculates virtual pool balances to return a price the pool is willing to offer
func (l *ConcentratedLiquidityPool) FairPrice() float64 {

	// pool is long
	if l.position > 0.0 {
		// x_v = P + (L / sqrt(pl))
		vX := l.position + (l.lowerL / math.Sqrt(l.base))

		// y_v = cc * rf + (L / sqrt(pl))
		cc := 0.0 // balance across all pool accounts???
		vY := cc*l.rf + (l.lower / math.Sqrt(l.lower))

		return vY / vX

	}

	// pool is short
	if l.position < 0.0 {
		// x_v = P + (cc * rf / pu) + (L / sqrt(pl)), pl is base price, pu upper price
		cc := 0.0 // balance across all pool accounts???
		vX := l.position + (cc * l.rf / l.upper) + (l.upperL / math.Sqrt((l.base)))

		// v_y = abs(P)*p_e + L*pl where pe average-entry price of position, pl is base price
		ae := 0.0
		vY := (-l.position)*ae + (l.upperL * l.base)
		return vY / vX

	}

	// position is 0
	return l.base
}

func (l *ConcentratedLiquidityPool) VolumeBetweenPrices(st, nd float64) float64 {

	return 0.0

}

// TODO track balance, track position, track average-entry-price

// buying and selling positions on the market in exchange for quote asset
// as price lowers pool is buying positions reducing "currency"
// as price rises pool is selling increasing "currency"
