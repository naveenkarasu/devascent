func pacific_atlantic(heights [][]int) [][]int {
	if len(heights) == 0 {
		return nil
	}
	m, n := len(heights), len(heights[0])
	type point struct{ i, j int }

	bfs := func(starts []point) map[point]bool {
		visited := make(map[point]bool)
		queue := make([]point, len(starts))
		copy(queue, starts)
		for _, s := range starts {
			visited[s] = true
		}
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, d := range dirs {
				ni, nj := cur.i+d[0], cur.j+d[1]
				np := point{ni, nj}
				if ni >= 0 && ni < m && nj >= 0 && nj < n && !visited[np] && heights[ni][nj] >= heights[cur.i][cur.j] {
					visited[np] = true
					queue = append(queue, np)
				}
			}
		}
		return visited
	}

	var pacStarts, atlStarts []point
	for j := 0; j < n; j++ {
		pacStarts = append(pacStarts, point{0, j})
		atlStarts = append(atlStarts, point{m - 1, j})
	}
	for i := 1; i < m; i++ {
		pacStarts = append(pacStarts, point{i, 0})
	}
	for i := 0; i < m-1; i++ {
		atlStarts = append(atlStarts, point{i, n - 1})
	}

	pac := bfs(pacStarts)
	atl := bfs(atlStarts)

	var result [][]int
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			p := point{i, j}
			if pac[p] && atl[p] {
				result = append(result, []int{i, j})
			}
		}
	}
	return result
}
