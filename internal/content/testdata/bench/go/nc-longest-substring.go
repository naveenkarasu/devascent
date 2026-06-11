func length_of_longest_substring(s string) int {
	lastSeen := map[byte]int{}
	left := 0
	best := 0
	for right := 0; right < len(s); right++ {
		ch := s[right]
		if idx, ok := lastSeen[ch]; ok && idx >= left {
			left = idx + 1
		}
		lastSeen[ch] = right
		if right-left+1 > best {
			best = right - left + 1
		}
	}
	return best
}
