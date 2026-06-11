func gcd_lcm(a, b int) []int {
	gcdFn := func(x, y int) int {
		for y != 0 {
			x, y = y, x%y
		}
		return x
	}
	g := gcdFn(a, b)
	l := a * b / g
	return []int{g, l}
}
