func int_to_roman(num int) string {
	mapping := [][2]interface{}{
		{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
		{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
		{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
	}
	result := []byte{}
	for _, pair := range mapping {
		value := pair[0].(int)
		numeral := pair[1].(string)
		for num >= value {
			result = append(result, numeral...)
			num -= value
		}
	}
	return string(result)
}
