func min_window(s string, t string) string {
	if len(t) == 0 || len(s) == 0 {
		return ""
	}
	need := map[byte]int{}
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	have := map[byte]int{}
	formed := 0
	required := len(need)
	left := 0
	bestLen := len(s) + 1
	bestStart := 0
	for right := 0; right < len(s); right++ {
		ch := s[right]
		have[ch]++
		if need[ch] > 0 && have[ch] == need[ch] {
			formed++
		}
		for formed == required {
			if right-left+1 < bestLen {
				bestLen = right - left + 1
				bestStart = left
			}
			leftCh := s[left]
			have[leftCh]--
			if need[leftCh] > 0 && have[leftCh] < need[leftCh] {
				formed--
			}
			left++
		}
	}
	if bestLen == len(s)+1 {
		return ""
	}
	return s[bestStart : bestStart+bestLen]
}
