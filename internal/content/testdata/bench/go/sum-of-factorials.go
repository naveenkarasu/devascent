func sum_of_factorials(n int) int {
	factorial := func(x int) int {
		result := 1
		for i := 2; i <= x; i++ {
			result *= i
		}
		return result
	}
	total := 0
	for i := 1; i <= n; i++ {
		total += factorial(i)
	}
	return total
}
