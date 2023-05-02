package pool

import (
	"errors"
	"fmt"
)

var (
	ErrNotEnoughBaseAsset  = errors.New("not enough base asset")
	ErrNotEnoughQuoteAsset = errors.New("not enough quote asset")
)

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

func NewConstantProductPool(x, y, k uint64) *CPMM {
	return &CPMM{
		k:     k,
		base:  x,
		quote: y,
	}
}

// GetTrade return a potential trade against this pool give the price and size
func (p *CPMM) Trade(size uint64, side Side) (uint64, error) {
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
	fmt.Println(p.base, p.quote, p.k)

	var x, y, dy uint64
	switch side {
	case SideBuy:
		x, y = p.base, p.quote
		dy = (y * size) / (x + size)
		fmt.Println("want to buy", size, "of asset X in exchange for", dy, "of Y")

		if size > p.base {
			return 0, ErrNotEnoughBaseAsset
		}
		p.base -= size
		p.quote += dy
	case SideSell:
		x, y = p.quote, p.base
		dy = (y * size) / (x + size)
		fmt.Println("want to sell", size, "of asset X in exchange for", dy, "of Y")
		if dy > p.quote {
			return 0, ErrNotEnoughQuoteAsset
		}
		p.base += size
		p.quote -= dy
	default:
		panic("unknown side")
	}
	return dy, nil
}

// Verify dubg function that just makes sure everything is ok in the pool
func (p *CPMM) Verify() error {
	if (p.base * p.quote) != p.k {
		return fmt.Errorf("pool not constant %d %d %d", p.base, p.quote, p.k)
	}

	return nil
}
