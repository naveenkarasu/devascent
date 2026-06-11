func max_product(nums []int) int {
	maxProd, minProd, res := nums[0], nums[0], nums[0]
	for _, n := range nums[1:] {
		a, b, c := n, maxProd*n, minProd*n
		maxProd = max(a, max(b, c))
		minProd = min(a, min(b, c))
		res = max(res, maxProd)
	}
	return res
}
