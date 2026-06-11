func lexical_order(n int) []int {
	result := []int{}
	cur := 1
	for len(result) < n {
		result = append(result, cur)
		if cur*10 <= n {
			cur *= 10
		} else {
			for cur%10 == 9 || cur >= n {
				cur /= 10
			}
			cur++
		}
	}
	return result
}
