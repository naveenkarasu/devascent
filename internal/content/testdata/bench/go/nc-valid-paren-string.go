func check_valid_string(s string) bool {
	lo, hi := 0, 0
	for _, c := range s {
		if c == '(' {
			lo++
			hi++
		} else if c == ')' {
			lo--
			hi--
		} else {
			lo--
			hi++
		}
		if hi < 0 {
			return false
		}
		if lo < 0 {
			lo = 0
		}
	}
	return lo == 0
}
