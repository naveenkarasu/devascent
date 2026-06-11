import "container/heap"

type swimEntry struct {
	t, i, j int
}

type swimHeap []swimEntry

func (h swimHeap) Len() int           { return len(h) }
func (h swimHeap) Less(a, b int) bool { return h[a].t < h[b].t }
func (h swimHeap) Swap(a, b int)      { h[a], h[b] = h[b], h[a] }
func (h *swimHeap) Push(x interface{}) { *h = append(*h, x.(swimEntry)) }
func (h *swimHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func swim_in_water(grid [][]int) int {
	n := len(grid)
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	h := &swimHeap{{grid[0][0], 0, 0}}
	heap.Init(h)
	visited[0][0] = true
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for h.Len() > 0 {
		e := heap.Pop(h).(swimEntry)
		if e.i == n-1 && e.j == n-1 {
			return e.t
		}
		for _, d := range dirs {
			ni, nj := e.i+d[0], e.j+d[1]
			if ni >= 0 && ni < n && nj >= 0 && nj < n && !visited[ni][nj] {
				visited[ni][nj] = true
				t := e.t
				if grid[ni][nj] > t {
					t = grid[ni][nj]
				}
				heap.Push(h, swimEntry{t, ni, nj})
			}
		}
	}
	return -1
}
