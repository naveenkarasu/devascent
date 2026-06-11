func filter_even_digit_numbers(nums []int) []int {
	digitCount := func(n int) int {
		if n == 0 {
			return 1
		}
		count := 0
		for n > 0 {
			count++
			n /= 10
		}
		return count
	}
	result := []int{}
	for _, x := range nums {
		if digitCount(x)%2 == 0 {
			result = append(result, x)
		}
	}
	return result
}
