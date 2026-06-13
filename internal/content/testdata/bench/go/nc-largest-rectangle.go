func largest_rectangle_area(heights []int) int {
	maxArea := 0
	type pair struct{ start, height int }
	stack := []pair{}
	for i, h := range heights {
		start := i
		for len(stack) > 0 && stack[len(stack)-1].height > h {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			area := top.height * (i - top.start)
			if area > maxArea {
				maxArea = area
			}
			start = top.start
		}
		stack = append(stack, pair{start, h})
	}
	for _, p := range stack {
		area := p.height * (len(heights) - p.start)
		if area > maxArea {
			maxArea = area
		}
	}
	return maxArea
}
