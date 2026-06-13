func sort_three_values(arr []int, low_val int, mid_val int, high_val int) []int {
	lo, mid, hi := 0, 0, len(arr)-1
	for mid <= hi {
		if arr[mid] == low_val {
			arr[lo], arr[mid] = arr[mid], arr[lo]
			lo++
			mid++
		} else if arr[mid] == mid_val {
			mid++
		} else {
			arr[mid], arr[hi] = arr[hi], arr[mid]
			hi--
		}
	}
	return arr
}
