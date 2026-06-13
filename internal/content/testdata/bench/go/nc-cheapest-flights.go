import "math"

func find_cheapest_price(n int, flights [][]int, src int, dst int, k int) int {
	prices := make([]int, n)
	for i := range prices {
		prices[i] = math.MaxInt
	}
	prices[src] = 0
	for round := 0; round <= k; round++ {
		tmp := make([]int, n)
		copy(tmp, prices)
		for _, f := range flights {
			u, v, w := f[0], f[1], f[2]
			if prices[u] != math.MaxInt && prices[u]+w < tmp[v] {
				tmp[v] = prices[u] + w
			}
		}
		prices = tmp
	}
	if prices[dst] == math.MaxInt {
		return -1
	}
	return prices[dst]
}
