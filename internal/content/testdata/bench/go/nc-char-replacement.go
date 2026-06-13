func character_replacement(s string, k int) int {
	counts := [26]int{}
	left := 0
	maxCount := 0
	best := 0
	for right := 0; right < len(s); right++ {
		ch := s[right] - 'A'
		counts[ch]++
		if counts[ch] > maxCount {
			maxCount = counts[ch]
		}
		for (right-left+1)-maxCount > k {
			counts[s[left]-'A']--
			left++
		}
		if right-left+1 > best {
			best = right - left + 1
		}
	}
	return best
}
