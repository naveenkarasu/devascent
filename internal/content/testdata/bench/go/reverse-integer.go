func reverse_integer(x int) int {
	sign := 1
	if x < 0 {
		sign = -1
		x = -x
	}
	rev := 0
	for x != 0 {
		rev = rev*10 + x%10
		x /= 10
	}
	rev *= sign
	if rev > (1<<31)-1 || rev < -(1<<31) {
		return 0
	}
	return rev
}
