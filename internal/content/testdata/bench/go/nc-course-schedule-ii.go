import "container/heap"

type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func find_order(num_courses int, prerequisites [][]int) []int {
	inDegree := make([]int, num_courses)
	adj := make([][]int, num_courses)
	for _, pre := range prerequisites {
		a, b := pre[0], pre[1]
		adj[b] = append(adj[b], a)
		inDegree[a]++
	}
	h := &intHeap{}
	for i := 0; i < num_courses; i++ {
		if inDegree[i] == 0 {
			heap.Push(h, i)
		}
	}
	heap.Init(h)
	order := []int{}
	for h.Len() > 0 {
		node := heap.Pop(h).(int)
		order = append(order, node)
		for _, nei := range adj[node] {
			inDegree[nei]--
			if inDegree[nei] == 0 {
				heap.Push(h, nei)
			}
		}
	}
	if len(order) == num_courses {
		return order
	}
	return []int{}
}
