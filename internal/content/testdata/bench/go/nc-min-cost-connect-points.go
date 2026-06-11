import "container/heap"

type minCostEntry struct {
	cost int
	idx  int
}

type minCostHeap []minCostEntry

func (h minCostHeap) Len() int            { return len(h) }
func (h minCostHeap) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h minCostHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minCostHeap) Push(x interface{}) { *h = append(*h, x.(minCostEntry)) }
func (h *minCostHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func abs_int(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min_cost_connect_points(points [][]int) int {
	n := len(points)
	visited := make([]bool, n)
	h := &minCostHeap{{0, 0}}
	heap.Init(h)
	total := 0
	count := 0
	for h.Len() > 0 && count < n {
		e := heap.Pop(h).(minCostEntry)
		if visited[e.idx] {
			continue
		}
		visited[e.idx] = true
		total += e.cost
		count++
		xi, yi := points[e.idx][0], points[e.idx][1]
		for j := 0; j < n; j++ {
			if !visited[j] {
				dist := abs_int(xi-points[j][0]) + abs_int(yi-points[j][1])
				heap.Push(h, minCostEntry{dist, j})
			}
		}
	}
	return total
}
