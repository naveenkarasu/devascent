func max_profit(prices []int) int {
	hold := -1 << 62
	sold := 0
	rest := 0
	for _, price := range prices {
		prevSold := sold
		sold = hold + price
		if rest-price > hold {
			hold = rest - price
		}
		if prevSold > rest {
			rest = prevSold
		}
	}
	if sold > rest {
		return sold
	}
	return rest
}
