func compressed_length(s string) int {
	n := len(s)
	count := 0
	i := 0
	for i < n {
		if i+1 < n && s[i] == s[i+1] {
			count++
			i += 2
		} else {
			count++
			i++
		}
	}
	return count
}
