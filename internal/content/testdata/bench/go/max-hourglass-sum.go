func max_hourglass_sum(grid [][]int) int {
	best := -1<<63
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			s := grid[i][j] + grid[i][j+1] + grid[i][j+2] +
				grid[i+1][j+1] +
				grid[i+2][j] + grid[i+2][j+1] + grid[i+2][j+2]
			if s > best {
				best = s
			}
		}
	}
	return best
}
