func can_finish(num_courses int, prerequisites [][]int) bool {
	graph := make([][]int, num_courses)
	for i := range graph {
		graph[i] = []int{}
	}
	for _, p := range prerequisites {
		graph[p[0]] = append(graph[p[0]], p[1])
	}
	// state: 0=unvisited, 1=visiting, 2=done
	state := make([]int, num_courses)
	var dfs func(c int) bool
	dfs = func(c int) bool {
		if state[c] == 1 {
			return false
		}
		if state[c] == 2 {
			return true
		}
		state[c] = 1
		for _, nxt := range graph[c] {
			if !dfs(nxt) {
				return false
			}
		}
		state[c] = 2
		return true
	}
	for c := 0; c < num_courses; c++ {
		if !dfs(c) {
			return false
		}
	}
	return true
}
