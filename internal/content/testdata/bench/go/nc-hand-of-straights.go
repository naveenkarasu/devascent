import "sort"

func is_n_straight_hand(hand []int, group_size int) bool {
	if len(hand)%group_size != 0 {
		return false
	}
	count := make(map[int]int)
	for _, c := range hand {
		count[c]++
	}
	keys := make([]int, 0, len(count))
	for k := range count {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, card := range keys {
		if count[card] > 0 {
			need := count[card]
			for i := 0; i < group_size; i++ {
				if count[card+i] < need {
					return false
				}
				count[card+i] -= need
			}
		}
	}
	return true
}
