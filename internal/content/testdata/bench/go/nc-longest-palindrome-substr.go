func longest_palindrome(s string) string {
	resStart, resLen := 0, 1

	expand := func(left, right int) {
		for left >= 0 && right < len(s) && s[left] == s[right] {
			length := right - left + 1
			if length > resLen {
				resStart, resLen = left, length
			}
			left--
			right++
		}
	}

	for i := range s {
		expand(i, i)
		expand(i, i+1)
	}

	return s[resStart : resStart+resLen]
}
