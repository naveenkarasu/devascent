import (
	"container/heap"
	"sort"
)

type intervalHeap [][2]int // [size, end]

func (h intervalHeap) Len() int            { return len(h) }
func (h intervalHeap) Less(i, j int) bool  { return h[i][0] < h[j][0] }
func (h intervalHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intervalHeap) Push(x interface{}) { *h = append(*h, x.([2]int)) }
func (h *intervalHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func min_interval(intervals [][]int, queries []int) []int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	n := len(queries)
	indexed := make([][2]int, n)
	for i, q := range queries {
		indexed[i] = [2]int{q, i}
	}
	sort.Slice(indexed, func(i, j int) bool {
		return indexed[i][0] < indexed[j][0]
	})
	res := make([]int, n)
	for i := range res {
		res[i] = -1
	}
	h := &intervalHeap{}
	heap.Init(h)
	idx := 0
	for _, qi := range indexed {
		q, origIdx := qi[0], qi[1]
		for idx < len(intervals) && intervals[idx][0] <= q {
			l, r := intervals[idx][0], intervals[idx][1]
			heap.Push(h, [2]int{r - l + 1, r})
			idx++
		}
		for h.Len() > 0 && (*h)[0][1] < q {
			heap.Pop(h)
		}
		if h.Len() > 0 {
			res[origIdx] = (*h)[0][0]
		}
	}
	return res
}
