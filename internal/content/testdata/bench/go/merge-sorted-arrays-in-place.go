func merge_sorted_in_place(arr1, arr2 []int) []int {
	m := len(arr1)
	n := len(arr2)
	result := make([]int, m+n)
	l, r, k := m-1, n-1, m+n-1
	for l >= 0 && r >= 0 {
		if arr1[l] > arr2[r] {
			result[k] = arr1[l]
			l--
		} else {
			result[k] = arr2[r]
			r--
		}
		k--
	}
	for l >= 0 {
		result[k] = arr1[l]
		l--
		k--
	}
	for r >= 0 {
		result[k] = arr2[r]
		r--
		k--
	}
	return result
}
