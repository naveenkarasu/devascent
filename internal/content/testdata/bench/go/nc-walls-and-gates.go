func walls_and_gates(rooms [][]int) [][]int {
	if len(rooms) == 0 || len(rooms[0]) == 0 {
		return rooms
	}
	const INF = 2147483647
	m, n := len(rooms), len(rooms[0])
	type point struct{ r, c int }
	queue := []point{}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if rooms[i][j] == 0 {
				queue = append(queue, point{i, j})
			}
		}
	}
	dirs := []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			nr, nc := cur.r+d.r, cur.c+d.c
			if nr >= 0 && nr < m && nc >= 0 && nc < n && rooms[nr][nc] == INF {
				rooms[nr][nc] = rooms[cur.r][cur.c] + 1
				queue = append(queue, point{nr, nc})
			}
		}
	}
	return rooms
}
