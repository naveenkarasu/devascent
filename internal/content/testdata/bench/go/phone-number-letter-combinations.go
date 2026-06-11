func letter_combinations(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}
	mapping := map[byte]string{
		'2': "abc", '3': "def", '4': "ghi", '5': "jkl",
		'6': "mno", '7': "pqrs", '8': "tuv", '9': "wxyz",
	}
	results := []string{""}
	for i := 0; i < len(digits); i++ {
		d := digits[i]
		letters := mapping[d]
		next := []string{}
		for _, prev := range results {
			for _, ch := range letters {
				next = append(next, prev+string(ch))
			}
		}
		results = next
	}
	return results
}
