package main

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
}
