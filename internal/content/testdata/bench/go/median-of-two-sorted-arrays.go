func find_median_sorted_arrays(nums1, nums2 []int) float64 {
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	m, n := len(nums1), len(nums2)
	total := m + n
	lo, hi := 0, m
	for lo <= hi {
		cut1 := (lo + hi) / 2
		cut2 := (total+1)/2 - cut1
		var l1, l2, r1, r2 float64
		if cut1 > 0 {
			l1 = float64(nums1[cut1-1])
		} else {
			l1 = -1e18
		}
		if cut2 > 0 {
			l2 = float64(nums2[cut2-1])
		} else {
			l2 = -1e18
		}
		if cut1 < m {
			r1 = float64(nums1[cut1])
		} else {
			r1 = 1e18
		}
		if cut2 < n {
			r2 = float64(nums2[cut2])
		} else {
			r2 = 1e18
		}
		if l1 <= r2 && l2 <= r1 {
			if total%2 == 0 {
				max1 := l1
				if l2 > max1 {
					max1 = l2
				}
				min1 := r1
				if r2 < min1 {
					min1 = r2
				}
				return (max1 + min1) / 2.0
			}
			if l1 > l2 {
				return l1
			}
			return l2
		} else if l1 > r2 {
			hi = cut1 - 1
		} else {
			lo = cut1 + 1
		}
	}
	return 0.0
}
