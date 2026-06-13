func nth_prime(n int) int {
	primes := []int{}
	i := 2
	for len(primes) < n {
		isPrime := true
		for _, p := range primes {
			if i%p == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			primes = append(primes, i)
		}
		i++
	}
	return primes[len(primes)-1]
}
