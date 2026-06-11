import "sort"

func letter_combinations(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}
	mapping := map[byte]string{
		'2': "abc", '3': "def", '4': "ghi", '5': "jkl",
		'6': "mno", '7': "pqrs", '8': "tuv", '9': "wxyz",
	}
	res := []string{}
	var backtrack func(index int, current string)
	backtrack = func(index int, current string) {
		if index == len(digits) {
			res = append(res, current)
			return
		}
		for _, ch := range mapping[digits[index]] {
			backtrack(index+1, current+string(ch))
		}
	}
	backtrack(0, "")
	sort.Strings(res)
	return res
}
