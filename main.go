package main

import (
	"fmt"

	"github.com/wwestgarth/liq-pool-amm/pool"
)

type Pool interface {
	// Return the best price for the given trade to buy/sell size of the base asset
	BestPrice(size float64, side pool.Side) (float64, error)

	// Actually do the trade
	Trade(size float64, side pool.Side) (float64, error)
}

type Pools struct {
	pools []Pool
}

func (ps *Pools) Add(pool Pool) {
	ps.pools = append(ps.pools, pool)
}

func (ps *Pools) Trade(price float64, size float64, side pool.Side) {
	var bestPool Pool
	var bestPrice float64
	found := false
	for _, p := range ps.pools {

		price, err := p.BestPrice(size, side)
		if err != nil {
			continue
		}

		if !found {
			found = true
			bestPool = p
			bestPrice = price
			continue
		}
		switch side {
		case pool.SideBuy:
			if price < bestPrice {
				bestPrice = price
				bestPool = p
			}
		case pool.SideSell:
			if price > bestPrice {
				bestPrice = price
				bestPool = p
			}
		}
	}

	if !found {
		// no trade, it will just sit on the normal order book
		return
	}

	dy, err := bestPool.Trade(size, side)
	if err != nil {
		panic("pool gave a best price but then couldn't trade it, should never happen")
	}
	fmt.Println("you traded", size, "of the base asset for", dy, "of the quote asset")
}

func main() {
	// Thoughts:
	// in vega when we ask the pool for a swap, we generate a trade. Does the sold asset go back into the pool balancing it, or does it go into
	// the LPs account?

	// Don't forget "price" is the amount of quote asset, so if I have an order for to sell 1000 base asset with a price of 500 quote asset,
	// I will ask the ask a pool "I want to give you (sell) 1000 base asset, how many quote asset will you give me?"
	// if they respond with N, the the orders cross is N >= 500 i.e the pool will give me more (a better price) that in my order.
	//
	// Flipped, if I have an order to buy 1000 base asset with a price of 500 quote asset, I will ask the pool
	// ""I want to take from you (buy) 1000 base asset, how many quote asset will I need to give you?"
	// if they response with N, then the orders cross if N <= 500 i. I have to give the pool fewer than I wanted to

	// make some pools

	pools := Pools{}

	pools.Add(pool.NewConstantProductPool(100, 100, 100*100))
	pools.Add(pool.NewConstantProductPool(125, 100, 125*100))

	// I want to sell 10 of the base asset for 9 of the quote asset
	pools.Trade(9, 10, pool.SideSell)
}
