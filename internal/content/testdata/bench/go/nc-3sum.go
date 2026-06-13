import "sort"

func three_sum(nums []int) [][]int {
	sort.Ints(nums)
	res := [][]int{}
	n := len(nums)
	for i := 0; i < n; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		lo, hi := i+1, n-1
		for lo < hi {
			total := nums[i] + nums[lo] + nums[hi]
			if total < 0 {
				lo++
			} else if total > 0 {
				hi--
			} else {
				res = append(res, []int{nums[i], nums[lo], nums[hi]})
				lo++
				hi--
				for lo < hi && nums[lo] == nums[lo-1] {
					lo++
				}
				for lo < hi && nums[hi] == nums[hi+1] {
					hi--
				}
			}
		}
	}
	return res
}
