func toggle_case(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch >= 'A' && ch <= 'Z' {
			result[i] = ch + 32
		} else if ch >= 'a' && ch <= 'z' {
			result[i] = ch - 32
		} else {
			result[i] = ch
		}
	}
	return string(result)
}
