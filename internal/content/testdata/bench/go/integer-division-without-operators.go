func divide_integers(dividend int, divisor int) int {
	const INT_MAX = 1<<31 - 1
	const INT_MIN = -(1 << 31)
	if dividend == INT_MIN && divisor == -1 {
		return INT_MAX
	}
	negative := (dividend < 0) != (divisor < 0)
	a, b := dividend, divisor
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	result := 0
	for a >= b {
		temp, multiple := b, 1
		for a >= (temp << 1) {
			temp <<= 1
			multiple <<= 1
		}
		a -= temp
		result += multiple
	}
	if negative {
		return -result
	}
	return result
}
