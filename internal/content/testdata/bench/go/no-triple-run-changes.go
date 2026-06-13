func min_changes_no_triple(s string) int {
	chars := []byte(s)
	ans := 0
	same := 1
	for i := 1; i < len(chars); i++ {
		if chars[i] == chars[i-1] {
			same++
		} else {
			same = 1
		}
		if same == 3 {
			ans++
			chars[i] = '@'
			same = 1
		}
	}
	return ans
}
