func get_sum(a, b int) int {
	mask := 0xFFFFFFFF
	for b&mask != 0 {
		carry := ((a & b) << 1) & mask
		a = (a ^ b) & mask
		b = carry
	}
	a &= mask
	if a <= 0x7FFFFFFF {
		return a
	}
	return ^(a ^ mask)
}
