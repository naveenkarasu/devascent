func steps_to_overtake(a int, b int) int {
	steps := 0
	for a <= b {
		a *= 3
		b *= 2
		steps++
	}
	return steps
}
