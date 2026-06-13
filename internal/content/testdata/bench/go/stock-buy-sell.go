func max_profit(prices []int) int {
	buy := prices[0]
	res := 0
	for i := 1; i < len(prices); i++ {
		profit := prices[i] - buy
		if profit > res {
			res = profit
		}
		if prices[i] < buy {
			buy = prices[i]
		}
	}
	return res
}
