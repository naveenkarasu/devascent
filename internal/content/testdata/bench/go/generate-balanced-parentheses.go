func generate_parentheses(n int) []string {
	results := []string{}
	var backtrack func(s string, open, close int)
	backtrack = func(s string, open, close int) {
		if len(s) == 2*n {
			results = append(results, s)
			return
		}
		if open < n {
			backtrack(s+"(", open+1, close)
		}
		if close < open {
			backtrack(s+")", open, close+1)
		}
	}
	backtrack("", 0, 0)
	return results
}
