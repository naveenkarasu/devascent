func split_odd_even(s string) []string {
	evenChars := []byte{}
	oddChars := []byte{}
	for i := 0; i < len(s); i++ {
		if i%2 == 0 {
			evenChars = append(evenChars, s[i])
		} else {
			oddChars = append(oddChars, s[i])
		}
	}
	return []string{string(evenChars), string(oddChars)}
}
