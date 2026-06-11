import "container/heap"

type distEntry struct {
	dist, x, y int
}

type distHeap []distEntry

func (h distHeap) Len() int            { return len(h) }
func (h distHeap) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h distHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *distHeap) Push(x interface{}) { *h = append(*h, x.(distEntry)) }
func (h *distHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func k_closest(points [][]int, k int) [][]int {
	h := &distHeap{}
	heap.Init(h)
	for _, p := range points {
		x, y := p[0], p[1]
		heap.Push(h, distEntry{x*x + y*y, x, y})
	}
	result := make([][]int, 0, k)
	for i := 0; i < k; i++ {
		e := heap.Pop(h).(distEntry)
		result = append(result, []int{e.x, e.y})
	}
	return result
}
