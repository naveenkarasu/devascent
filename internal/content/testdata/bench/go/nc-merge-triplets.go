func merge_triplets(triplets [][]int, target []int) bool {
	result := []int{0, 0, 0}
	for _, t := range triplets {
		if t[0] <= target[0] && t[1] <= target[1] && t[2] <= target[2] {
			for i := 0; i < 3; i++ {
				if t[i] > result[i] {
					result[i] = t[i]
				}
			}
		}
	}
	return result[0] == target[0] && result[1] == target[1] && result[2] == target[2]
}
