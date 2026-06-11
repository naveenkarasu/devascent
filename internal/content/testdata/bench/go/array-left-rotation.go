func left_rotate(arr []int, k int) []int {
	n := len(arr)
	if n == 0 {
		return arr
	}
	k = k % n
	result := make([]int, n)
	copy(result, arr[k:])
	copy(result[n-k:], arr[:k])
	return result
}
