func check_inclusion(s1 string, s2 string) bool {
	if len(s1) > len(s2) {
		return false
	}
	need := [26]int{}
	window := [26]int{}
	for i := 0; i < len(s1); i++ {
		need[s1[i]-'a']++
	}
	k := len(s1)
	for i := 0; i < len(s2); i++ {
		window[s2[i]-'a']++
		if i >= k {
			window[s2[i-k]-'a']--
		}
		if window == need {
			return true
		}
	}
	return false
}
