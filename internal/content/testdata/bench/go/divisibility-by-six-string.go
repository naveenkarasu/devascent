func is_divisible_by_six(s string) bool {
	if len(s) == 0 {
		return false
	}
	lastDigit := int(s[len(s)-1] - '0')
	if lastDigit%2 != 0 {
		return false
	}
	digitSum := 0
	for _, ch := range s {
		digitSum += int(ch - '0')
	}
	return digitSum%3 == 0
}
