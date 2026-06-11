func count_vowels(s string) int {
	vowels := map[rune]bool{
		'a': true, 'e': true, 'i': true, 'o': true, 'u': true,
		'A': true, 'E': true, 'I': true, 'O': true, 'U': true,
	}
	count := 0
	for _, ch := range s {
		if vowels[ch] {
			count++
		}
	}
	return count
}
