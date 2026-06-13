func can_jump(nums []int) bool {
	reach := 0
	for i, n := range nums {
		if i > reach {
			return false
		}
		if i+n > reach {
			reach = i + n
		}
	}
	return true
}
