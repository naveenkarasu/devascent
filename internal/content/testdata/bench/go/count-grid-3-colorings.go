func count_grid_colorings(n int, m int) int {
	const MOD = 1_000_000_007
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, m)
	}

	var backtrack func(pos int) int
	backtrack = func(pos int) int {
		if pos == n*m {
			return 1
		}
		r, c := pos/m, pos%m
		total := 0
		for color := 1; color <= 3; color++ {
			if (r == 0 || grid[r-1][c] != color) &&
				(c == 0 || grid[r][c-1] != color) {
				grid[r][c] = color
				total = (total + backtrack(pos+1)) % MOD
				grid[r][c] = 0
			}
		}
		return total
	}

	return backtrack(0)
}
