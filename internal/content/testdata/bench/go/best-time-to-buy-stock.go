func max_profit(prices []int) int {
	minPrice := prices[0]
	profit := 0
	for _, p := range prices[1:] {
		if p-minPrice > profit {
			profit = p - minPrice
		}
		if p < minPrice {
			minPrice = p
		}
	}
	return profit
}
