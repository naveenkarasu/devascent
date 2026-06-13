func longest_consecutive(nums []int) int {
	s := map[int]bool{}
	for _, n := range nums {
		s[n] = true
	}
	best := 0
	for n := range s {
		if !s[n-1] {
			length := 1
			for s[n+length] {
				length++
			}
			if length > best {
				best = length
			}
		}
	}
	return best
}
