func roman_to_int(s string) int {
	doubles := map[string]int{
		"CM": 900, "CD": 400, "XC": 90, "XL": 40, "IX": 9, "IV": 4,
	}
	singles := map[string]int{
		"M": 1000, "D": 500, "C": 100, "L": 50, "X": 10, "V": 5, "I": 1,
	}
	total := 0
	i := 0
	for i < len(s) {
		if i < len(s)-1 {
			two := string(s[i : i+2])
			if v, ok := doubles[two]; ok {
				total += v
				i += 2
				continue
			}
		}
		total += singles[string(s[i:i+1])]
		i++
	}
	return total
}
