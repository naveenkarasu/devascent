import "container/heap"

type maxHeap []int
func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type minHeap []int
func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type MedianFinder struct {
	lo *maxHeap // lower half, max-heap
	hi *minHeap // upper half, min-heap
}

func newMedianFinder() *MedianFinder {
	lo := &maxHeap{}
	hi := &minHeap{}
	heap.Init(lo)
	heap.Init(hi)
	return &MedianFinder{lo: lo, hi: hi}
}

func (m *MedianFinder) addNum(v int) {
	heap.Push(m.lo, v)
	heap.Push(m.hi, heap.Pop(m.lo).(int))
	if m.hi.Len() > m.lo.Len() {
		heap.Push(m.lo, heap.Pop(m.hi).(int))
	}
}

func (m *MedianFinder) findMedian() float64 {
	if m.lo.Len() > m.hi.Len() {
		return float64((*m.lo)[0])
	}
	return float64((*m.lo)[0]+(*m.hi)[0]) / 2.0
}

func median_stream_ops(operations [][]interface{}) []interface{} {
	mf := newMedianFinder()
	out := make([]interface{}, len(operations))
	for i, op := range operations {
		switch op[0].(string) {
		case "addNum":
			mf.addNum(op[1].(int))
			out[i] = nil
		default: // findMedian
			out[i] = mf.findMedian()
		}
	}
	return out
}
