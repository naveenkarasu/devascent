func two_sum_ii(numbers []int, target int) []int {
	i, j := 0, len(numbers)-1
	for i < j {
		total := numbers[i] + numbers[j]
		if total == target {
			return []int{i + 1, j + 1}
		} else if total < target {
			i++
		} else {
			j--
		}
	}
	return []int{}
}
