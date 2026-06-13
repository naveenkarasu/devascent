import "math"

func find_median_sorted_arrays(nums1 []int, nums2 []int) int {
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	m, n := len(nums1), len(nums2)
	half := (m + n) / 2
	lo, hi := 0, m
	for {
		i := (lo + hi) / 2
		j := half - i
		var left1, right1, left2, right2 int
		if i > 0 {
			left1 = nums1[i-1]
		} else {
			left1 = math.MinInt
		}
		if i < m {
			right1 = nums1[i]
		} else {
			right1 = math.MaxInt
		}
		if j > 0 {
			left2 = nums2[j-1]
		} else {
			left2 = math.MinInt
		}
		if j < n {
			right2 = nums2[j]
		} else {
			right2 = math.MaxInt
		}
		if left1 <= right2 && left2 <= right1 {
			if (m+n)%2 == 1 {
				if right1 < right2 {
					return right1
				}
				return right2
			}
			maxLeft := left1
			if left2 > maxLeft {
				maxLeft = left2
			}
			minRight := right1
			if right2 < minRight {
				minRight = right2
			}
			return (maxLeft + minRight) / 2
		} else if left1 > right2 {
			hi = i - 1
		} else {
			lo = i + 1
		}
	}
}
