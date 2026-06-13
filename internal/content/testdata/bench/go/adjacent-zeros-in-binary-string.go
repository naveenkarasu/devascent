func has_adjacent_zeros(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] == '0' && s[i-1] == '0' {
			return true
		}
	}
	return false
}
