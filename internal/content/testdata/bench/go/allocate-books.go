func allocate_books(pages []int, students int) int {
	n := len(pages)
	if students > n {
		return -1
	}
	canAllocate := func(maxPages int) bool {
		count := 1
		total := 0
		for _, p := range pages {
			if p > maxPages {
				return false
			}
			if total+p > maxPages {
				count++
				total = p
			} else {
				total += p
			}
		}
		return count <= students
	}
	lo := pages[0]
	hi := 0
	for _, p := range pages {
		if p > lo {
			lo = p
		}
		hi += p
	}
	ans := hi
	for lo <= hi {
		mid := (lo + hi) / 2
		if canAllocate(mid) {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return ans
}
