package pool

import (
	"fmt"
	"math"
)

type ConcentratedLiquidityPool struct {
	initial      float64 // initial commitment to the pool
	base         float64 // price where the pool has position 0, and contents it all physical tokens
	lower        float64 // price in [lower, base) pool will have a long position
	upper        float64 // price in (base, upper] pool will have a short position
	position     float64 // the pools current position
	balance      float64
	averageEntry float64

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
	fmt.Println(vw, "/ ( ", math.Sqrt(u), " - ", math.Sqrt(b), " ) = ", vw, " / ", math.Sqrt(u)-math.Sqrt(b))

	// lower so interval (base, lower]
	lowL := vw / (math.Sqrt(b) - math.Sqrt(l))

	fmt.Println("init L low:", lowL, ", L high:", upL)

	return &ConcentratedLiquidityPool{
		base:     b,
		lower:    l,
		upper:    u,
		initial:  i,
		balance:  i,
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
		fmt.Println("Xv = ", l.position, " + ( ", l.lowerL, " / ", "sqrt(", l.base, ") )")
		fmt.Println("Xv = ", l.position, " + ( ", l.lowerL, " / ", math.Sqrt(l.base), ")")

		// y_v = cc * rf + (L / sqrt(pl))

		// let
		// for small P,  y_v / x_v = (cc*rf)/ stuff + sqrt(base/lower)

		cc := l.balance // balance across all pool accounts???
		vY := cc*l.rf + (l.lowerL / math.Sqrt(l.lower))
		fmt.Println("Yv = ", cc, " + ", "( ", l.lowerL, " / sqrt(", l.lower, ") )")
		fmt.Println("Yv = ", cc, " + ", "( ", l.lowerL, " / ", math.Sqrt(l.lower), " )")
		fmt.Println("fair prices, x:", vX, "y:", vY, "fair price:", vY/vX)
		return vY / vX

	}

	// pool is short
	if l.position < 0.0 {
		// x_v = P + (cc * rf / pu) + (L / sqrt(pl)), pl is base price, pu upper price
		cc := l.balance // balance across all pool accounts???
		vX := l.position + (cc * l.rf / l.upper) + (l.upperL / math.Sqrt((l.base)))

		// v_y = abs(P)*p_e + L*pl where pe average-entry price of position, pl is base price
		ae := l.averageEntry
		fmt.Println("average entry", ae)
		vY := (-l.position)*ae + (l.upperL * l.base)
		return vY / vX

	}

	// position is 0
	return l.base
}

func (l *ConcentratedLiquidityPool) Trade(side Side, size float64) {

	if side == SideSell {
		// incoming order is selling so the pool is buying and becoming long
		price := l.FairPrice()
		sumproduct := l.averageEntry * l.position
		sumproduct += price * size

		l.position += size

		l.averageEntry = sumproduct / l.position
		return
	}

	if side == SideBuy {
		// incoming order is buying so the pool is selling and becoming long
		price := l.FairPrice()
		sumproduct := l.averageEntry * math.Abs(l.position)
		sumproduct += price * size

		l.position -= size

		l.averageEntry = sumproduct / math.Abs(l.position)
	}

}

func (l *ConcentratedLiquidityPool) VolumeBetweenPrices(st, nd float64) float64 {
	L := l.lowerL
	upperPrice := l.base
	if l.position < 0.0 {
		fmt.Println("short position")
		upperPrice = l.upper
		L = l.upperL
	}
	// TODO is position is zero, use the bounds of the curve we would move into

	fmt.Println("L:", L, "up:", upperPrice)
	sp := math.Sqrt(upperPrice)

	pp := math.Sqrt(st)
	stImplied := (L * (sp - pp)) / (sp * pp)

	pp = math.Sqrt(nd)
	ndImplied := (L * (sp - pp)) / (sp * pp)

	fmt.Println("impliedSt:", stImplied)
	fmt.Println("impliedNd:", ndImplied)
	return math.Abs(stImplied - ndImplied)

}

// TODO track balance, track position, track average-entry-price

// buying and selling positions on the market in exchange for quote asset
// as price lowers pool is buying positions reducing "currency"
// as price rises pool is selling increasing "currency"
