func daily_temperatures(temps []int) []int {
	result := make([]int, len(temps))
	type pair struct{ temp, idx int }
	stack := []pair{}
	for i, t := range temps {
		for len(stack) > 0 && t > stack[len(stack)-1].temp {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			result[top.idx] = i - top.idx
		}
		stack = append(stack, pair{t, i})
	}
	return result
}
