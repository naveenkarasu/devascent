func max_square_divisor_root(n int) int {
	answer := 1
	remaining := n
	for fact := 2; fact*fact <= n; fact++ {
		count := 0
		for remaining%fact == 0 {
			count++
			remaining /= fact
			if count == 2 {
				answer *= fact
				count = 0
			}
		}
	}
	return answer
}
