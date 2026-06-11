import "math"

func reverse_integer(x int) int {
	sign := 1
	if x < 0 {
		sign = -1
		x = -x
	}
	result := 0
	for x != 0 {
		result = result*10 + x%10
		x /= 10
	}
	result *= sign
	if result < math.MinInt32 || result > math.MaxInt32 {
		return 0
	}
	return result
}
