func is_perfect_square(num int) bool {
	if num < 1 {
		return false
	}
	lo, hi := 1, num
	for lo <= hi {
		mid := (lo + hi) / 2
		sq := mid * mid
		if sq == num {
			return true
		} else if sq < num {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return false
}
