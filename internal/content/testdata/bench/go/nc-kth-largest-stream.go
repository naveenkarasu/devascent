import "container/heap"

type kthMinHeap []int

func (h kthMinHeap) Len() int           { return len(h) }
func (h kthMinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h kthMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *kthMinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *kthMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type KthLargest struct {
	k    int
	heap *kthMinHeap
}

func newKthLargest(k int, nums []int) *KthLargest {
	h := &kthMinHeap{}
	heap.Init(h)
	kl := &KthLargest{k: k, heap: h}
	for _, v := range nums {
		heap.Push(h, v)
	}
	for h.Len() > k {
		heap.Pop(h)
	}
	return kl
}

func (kl *KthLargest) add(val int) int {
	heap.Push(kl.heap, val)
	for kl.heap.Len() > kl.k {
		heap.Pop(kl.heap)
	}
	return (*kl.heap)[0]
}

func kth_largest_stream_ops(k int, nums []int, operations [][]interface{}) []int {
	kl := newKthLargest(k, nums)
	out := make([]int, len(operations))
	for i, op := range operations {
		out[i] = kl.add(op[1].(int))
	}
	return out
}
