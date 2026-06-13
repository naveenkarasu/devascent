func find_majority_element(nums []int) int {
	candidate, count := 0, 0
	for _, x := range nums {
		if count == 0 {
			candidate = x
		}
		if x == candidate {
			count++
		} else {
			count--
		}
	}
	// Verify
	freq := 0
	for _, x := range nums {
		if x == candidate {
			freq++
		}
	}
	if freq > len(nums)/2 {
		return candidate
	}
	return -1
}
