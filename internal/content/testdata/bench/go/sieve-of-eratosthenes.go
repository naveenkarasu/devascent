func primes_up_to_n(n int) []int {
	if n < 2 {
		return []int{}
	}
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	result := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			result = append(result, i)
		}
	}
	return result
}
