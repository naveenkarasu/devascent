import "container/heap"

type delayEntry struct {
	dist int
	node int
}

type delayHeap []delayEntry

func (h delayHeap) Len() int            { return len(h) }
func (h delayHeap) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h delayHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *delayHeap) Push(x interface{}) { *h = append(*h, x.(delayEntry)) }
func (h *delayHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func network_delay_time(times [][]int, n int, k int) int {
	type edge struct{ to, w int }
	graph := make(map[int][]edge)
	for _, t := range times {
		u, v, w := t[0], t[1], t[2]
		graph[u] = append(graph[u], edge{v, w})
	}
	dist := make(map[int]int)
	h := &delayHeap{{0, k}}
	heap.Init(h)
	for h.Len() > 0 {
		e := heap.Pop(h).(delayEntry)
		if _, ok := dist[e.node]; ok {
			continue
		}
		dist[e.node] = e.dist
		for _, nb := range graph[e.node] {
			if _, ok := dist[nb.to]; !ok {
				heap.Push(h, delayEntry{e.dist + nb.w, nb.to})
			}
		}
	}
	if len(dist) != n {
		return -1
	}
	maxDist := 0
	for _, d := range dist {
		if d > maxDist {
			maxDist = d
		}
	}
	return maxDist
}
