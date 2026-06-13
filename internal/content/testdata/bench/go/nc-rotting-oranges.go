func oranges_rotting(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	type point struct{ i, j, t int }
	var queue []point
	fresh := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 2 {
				queue = append(queue, point{i, j, 0})
			} else if grid[i][j] == 1 {
				fresh++
			}
		}
	}
	if fresh == 0 {
		return 0
	}
	minutes := 0
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			ni, nj := cur.i+d[0], cur.j+d[1]
			if ni >= 0 && ni < m && nj >= 0 && nj < n && grid[ni][nj] == 1 {
				grid[ni][nj] = 2
				fresh--
				if cur.t+1 > minutes {
					minutes = cur.t + 1
				}
				queue = append(queue, point{ni, nj, cur.t + 1})
			}
		}
	}
	if fresh == 0 {
		return minutes
	}
	return -1
}
