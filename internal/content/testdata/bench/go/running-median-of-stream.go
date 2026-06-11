import "container/heap"

// maxHeap stores the lower half (inverted for max-heap behavior)
type maxH []int

func (h maxH) Len() int           { return len(h) }
func (h maxH) Less(i, j int) bool { return h[i] > h[j] }
func (h maxH) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxH) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *maxH) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// minH stores the upper half
type minH []int

func (h minH) Len() int           { return len(h) }
func (h minH) Less(i, j int) bool { return h[i] < h[j] }
func (h minH) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minH) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *minH) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func running_medians(stream []int) []float64 {
	lo := &maxH{}
	hi := &minH{}
	heap.Init(lo)
	heap.Init(hi)
	result := []float64{}
	for _, num := range stream {
		heap.Push(lo, num)
		heap.Push(hi, heap.Pop(lo))
		if hi.Len() > lo.Len() {
			heap.Push(lo, heap.Pop(hi))
		}
		if lo.Len() == hi.Len() {
			result = append(result, float64((*lo)[0]+(*hi)[0])/2.0)
		} else {
			result = append(result, float64((*lo)[0]))
		}
	}
	return result
}
