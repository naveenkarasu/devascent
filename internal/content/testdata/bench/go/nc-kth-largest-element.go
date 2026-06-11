import "container/heap"

type minHeapInt []int

func (h minHeapInt) Len() int            { return len(h) }
func (h minHeapInt) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeapInt) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeapInt) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *minHeapInt) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func find_kth_largest(nums []int, k int) int {
	h := &minHeapInt{}
	heap.Init(h)
	for _, n := range nums {
		heap.Push(h, n)
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	return (*h)[0]
}
