func selection_sort_k_passes(arr []int, k int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	n := len(result)
	passes := k
	if n-1 < passes {
		passes = n - 1
	}
	for i := 0; i < passes; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if result[j] < result[minIdx] {
				minIdx = j
			}
		}
		result[i], result[minIdx] = result[minIdx], result[i]
	}
	return result
}
