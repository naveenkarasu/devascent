func is_valid(s string) bool {
	stack := []rune{}
	pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
	for _, ch := range s {
		if open, ok := pairs[ch]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != open {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, ch)
		}
	}
	return len(stack) == 0
}
