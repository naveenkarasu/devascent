func merge_k_lists(lists [][]int) []int {
	// Use a simple k-way merge with a min-heap (implemented as slice)
	type entry struct {
		val  int
		list int
		idx  int
	}
	heap := []entry{}
	// heapify helpers
	var pushHeap func(e entry)
	var popHeap func() entry
	pushHeap = func(e entry) {
		heap = append(heap, e)
		i := len(heap) - 1
		for i > 0 {
			parent := (i - 1) / 2
			if heap[parent].val > heap[i].val {
				heap[parent], heap[i] = heap[i], heap[parent]
				i = parent
			} else {
				break
			}
		}
	}
	popHeap = func() entry {
		top := heap[0]
		last := len(heap) - 1
		heap[0] = heap[last]
		heap = heap[:last]
		i := 0
		for {
			left, right := 2*i+1, 2*i+2
			smallest := i
			if left < len(heap) && heap[left].val < heap[smallest].val {
				smallest = left
			}
			if right < len(heap) && heap[right].val < heap[smallest].val {
				smallest = right
			}
			if smallest == i {
				break
			}
			heap[i], heap[smallest] = heap[smallest], heap[i]
			i = smallest
		}
		return top
	}
	for i, lst := range lists {
		if len(lst) > 0 {
			pushHeap(entry{lst[0], i, 0})
		}
	}
	out := []int{}
	for len(heap) > 0 {
		e := popHeap()
		out = append(out, e.val)
		if e.idx+1 < len(lists[e.list]) {
			pushHeap(entry{lists[e.list][e.idx+1], e.list, e.idx + 1})
		}
	}
	return out
}
