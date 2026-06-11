import "math"

func perfect_numbers_up_to(n int) []int {
	result := []int{}
	for i := 2; i <= n; i++ {
		s := 1
		r := int(math.Sqrt(float64(i)))
		for d := 2; d <= r; d++ {
			if i%d == 0 {
				s += d + i/d
				if d*d == i {
					s -= d
				}
			}
		}
		if s == i {
			result = append(result, i)
		}
	}
	return result
}
