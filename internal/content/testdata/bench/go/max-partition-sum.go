func max_partition_sum(arr []int) int {
	total := 0
	for _, v := range arr {
		total += v
	}
	current := 0
	ans := 0
	for i := 0; i < len(arr)-1; i++ {
		current += arr[i]
		left := current
		right := total - current
		best := left
		if right > best {
			best = right
		}
		if best > ans {
			ans = best
		}
	}
	return ans
}
