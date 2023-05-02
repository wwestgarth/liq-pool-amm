package pool

import (
	"errors"
	"fmt"
	"math"
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
	k     float64
	base  float64 // the asset being purchased or sold X
	quote float64 // the asset which can be exchanged for the base asset Y
}

func NewConstantProductPool(x, y, k float64) *CPMM {
	return &CPMM{
		k:     k,
		base:  x,
		quote: y,
	}
}

// BestPrice given an order of size asset X being bought/sold return the change in asset Y
// i.e the amount that will be received, or will have to be supplied.
func (p *CPMM) BestPrice(size float64, side Side) (float64, error) {
	var x, y, dy float64
	switch side {
	case SideBuy:
		x, y = p.base, p.quote
		dy = (y * size) / (x + size)
		fmt.Println("want to buy", size, "of base asset, pool will need", dy, "of quote asset")

		if size > p.base {
			return 0, ErrNotEnoughBaseAsset
		}
	case SideSell:
		x, y = p.quote, p.base
		dy = (y * size) / (x + size)
		fmt.Println("want to sell", size, "of base asset, pool will give", dy, "of quote asset")
		if dy > p.quote {
			return 0, ErrNotEnoughQuoteAsset
		}
	default:
		panic("unknown side")
	}
	return dy, nil
}

// GetTrade return a potential trade against this pool give the price and size
func (p *CPMM) Trade(size float64, side Side) (float64, error) {
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

	var x, y, dy float64
	switch side {
	case SideBuy:
		x, y = p.base, p.quote
		dy = (y * size) / (x + size)
		fmt.Println("trading: receiving", size, "of base asset and giving", dy, "of quote asset")

		if size > p.base {
			return 0, ErrNotEnoughBaseAsset
		}
		p.base -= size
		p.quote += dy
	case SideSell:
		x, y = p.quote, p.base
		dy = (y * size) / (x + size)
		fmt.Println("trading: giving", size, "of asset base asset and receiving", dy, "of quote asset")
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

func (p *CPMM) Prices() (float64, float64) {
	return p.base / p.quote, p.quote / p.base
}

// Verify dubg function that just makes sure everything is ok in the pool
func (p *CPMM) Verify() error {
	// TODO: worry about rounding later
	tol := float64(0.000000001)
	diff := math.Abs((p.base * p.quote) - p.k)
	if diff > tol {
		fmt.Println(diff, tol)
		return fmt.Errorf("pool not constant %f %f %f", p.base, p.quote, p.k)
	}

	return nil
}
