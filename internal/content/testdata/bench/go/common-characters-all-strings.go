func count_common_characters(strings []string) int {
	if len(strings) == 0 {
		return 0
	}
	common := make(map[rune]bool)
	for _, ch := range strings[0] {
		common[ch] = true
	}
	for _, s := range strings[1:] {
		inS := make(map[rune]bool)
		for _, ch := range s {
			inS[ch] = true
		}
		for ch := range common {
			if !inS[ch] {
				delete(common, ch)
			}
		}
	}
	return len(common)
}
