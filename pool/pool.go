package pool

import "fmt"

type Side int

const ( // iota is reset to 0
	SideBuy Side = iota
	SideSell
)

type CPMM struct {
	// x * y = k
	k     uint64
	base  uint64 // the asset being purchased or sold X
	quote uint64 // the asset which can be exchanged for the base asset Y
}

func NewConstantProuctPool(x, y, k uint64) *CPMM {
	return &CPMM{
		k:     k,
		base:  x,
		quote: y,
	}
}

// GetTrade return a potential trade against this pool give the price and size
func (p *CPMM) Trade(size uint64, side Side) uint64 {
	// We have:
	// x * y = k

	// The constraint is:
	// (x + dx) * (y - dy) = k

	// which rearranges to:
	// dy = y - k / (x + dx)

	// we know k = x * y, so:
	// dy = y - (x * y) / (x + dx)

	// and eventually:
	// dy = (y * dx) / (x + dx)

	var x, y, dy uint64
	switch side {
	case SideBuy:
		x, y = p.base, p.quote
		dy = (y * size) / (x + size)
		p.base -= size
		p.quote += dy
		fmt.Println("want to buy", size, "of asset X in exchange for", dy, "of Y")
	case SideSell:
		fmt.Println("want to sell", size, "of asset X in exchange for Y")
		x, y = p.quote, p.base
		dy = (y * size) / (x + size)
		p.base += size
		p.quote -= dy
		fmt.Println("want to sell", size, "of asset X in exchange for", dy, "of Y")
	default:
		panic("unknown side")
	}
	return dy
}
