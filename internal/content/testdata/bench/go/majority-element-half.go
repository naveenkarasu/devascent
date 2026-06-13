func majority_element(nums []int) int {
	count := 0
	candidate := 0
	for _, n := range nums {
		if count == 0 {
			candidate = n
		}
		if n == candidate {
			count++
		} else {
			count--
		}
	}
	cnt := 0
	for _, n := range nums {
		if n == candidate {
			cnt++
		}
	}
	if cnt > len(nums)/2 {
		return candidate
	}
	return -1
}
