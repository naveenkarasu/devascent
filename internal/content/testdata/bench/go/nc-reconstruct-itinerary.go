import "sort"

func find_itinerary(tickets [][]string) []string {
	graph := make(map[string][]string)
	for _, t := range tickets {
		src, dst := t[0], t[1]
		graph[src] = append(graph[src], dst)
	}
	// Sort each adjacency list in reverse order so pop (from end) gives smallest first
	for src := range graph {
		sort.Sort(sort.Reverse(sort.StringSlice(graph[src])))
	}
	result := []string{}
	var dfs func(airport string)
	dfs = func(airport string) {
		for len(graph[airport]) > 0 {
			n := len(graph[airport])
			next := graph[airport][n-1]
			graph[airport] = graph[airport][:n-1]
			dfs(next)
		}
		result = append(result, airport)
	}
	dfs("JFK")
	// reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
