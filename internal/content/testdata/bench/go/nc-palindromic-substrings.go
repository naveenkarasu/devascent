func count_substrings(s string) int {
	count := 0

	expand := func(left, right int) {
		for left >= 0 && right < len(s) && s[left] == s[right] {
			count++
			left--
			right++
		}
	}

	for i := range s {
		expand(i, i)
		expand(i, i+1)
	}

	return count
}
