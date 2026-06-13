import (
	"math"
	"strconv"
)

func palindrome_primes_in_range(a, b int) []int {
	isPrime := func(n int) bool {
		if n < 2 {
			return false
		}
		for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
			if n%i == 0 {
				return false
			}
		}
		return true
	}
	isPalin := func(n int) bool {
		s := strconv.Itoa(n)
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			if s[i] != s[j] {
				return false
			}
		}
		return true
	}
	result := []int{}
	for i := a; i <= b; i++ {
		if isPrime(i) && isPalin(i) {
			result = append(result, i)
		}
	}
	return result
}
