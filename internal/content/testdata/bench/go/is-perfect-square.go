import "math"

func is_perfect_square(n int) bool {
	if n < 0 {
		return false
	}
	r := int(math.Sqrt(float64(n)))
	return r*r == n
}
