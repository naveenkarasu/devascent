func next_permutation(nums []int) []int {
	n := len(nums)
	pivot := -1
	for i := n - 2; i >= 0; i-- {
		if nums[i] < nums[i+1] {
			pivot = i
			break
		}
	}
	if pivot == -1 {
		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			nums[i], nums[j] = nums[j], nums[i]
		}
		return nums
	}
	for r := n - 1; r > pivot; r-- {
		if nums[r] > nums[pivot] {
			nums[pivot], nums[r] = nums[r], nums[pivot]
			break
		}
	}
	for i, j := pivot+1, n-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	return nums
}
