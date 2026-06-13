func min_max_pages(pages []int, num_students int) int {
	if num_students > len(pages) {
		return -1
	}
	canAllocate := func(maxPages int) bool {
		students := 1
		current := 0
		for _, p := range pages {
			if p > maxPages {
				return false
			}
			if current+p > maxPages {
				students++
				current = p
			} else {
				current += p
			}
		}
		return students <= num_students
	}
	lo := 0
	hi := 0
	for _, p := range pages {
		if p > lo {
			lo = p
		}
		hi += p
	}
	result := hi
	for lo <= hi {
		mid := (lo + hi) / 2
		if canAllocate(mid) {
			result = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return result
}
