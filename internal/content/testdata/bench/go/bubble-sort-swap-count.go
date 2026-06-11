func bubble_sort_info(arr []int) map[string]int {
	a := make([]int, len(arr))
	copy(a, arr)
	swaps := 0
	isSorted := false
	for !isSorted {
		isSorted = true
		for i := 0; i < len(a)-1; i++ {
			if a[i] > a[i+1] {
				a[i], a[i+1] = a[i+1], a[i]
				swaps++
				isSorted = false
			}
		}
	}
	return map[string]int{"swaps": swaps, "first": a[0], "last": a[len(a)-1]}
}
