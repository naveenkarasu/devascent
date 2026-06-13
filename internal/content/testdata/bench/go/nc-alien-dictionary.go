import (
	"container/heap"
	"sort"
)

type charHeap []byte

func (h charHeap) Len() int            { return len(h) }
func (h charHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h charHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *charHeap) Push(x interface{}) { *h = append(*h, x.(byte)) }
func (h *charHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func alien_order(words []string) string {
	// collect all chars
	indeg := make(map[byte]int)
	adj := make(map[byte][]byte)
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			if _, ok := indeg[w[i]]; !ok {
				indeg[w[i]] = 0
				adj[w[i]] = []byte{}
			}
		}
	}
	for i := 0; i < len(words)-1; i++ {
		a, b := words[i], words[i+1]
		m := len(a)
		if len(b) < m {
			m = len(b)
		}
		if len(a) > len(b) {
			eq := true
			for k := 0; k < m; k++ {
				if a[k] != b[k] {
					eq = false
					break
				}
			}
			if eq {
				return ""
			}
		}
		for j := 0; j < m; j++ {
			if a[j] != b[j] {
				// check if already added
				found := false
				for _, nb := range adj[a[j]] {
					if nb == b[j] {
						found = true
						break
					}
				}
				if !found {
					adj[a[j]] = append(adj[a[j]], b[j])
					indeg[b[j]]++
				}
				break
			}
		}
	}
	h := &charHeap{}
	for c, deg := range indeg {
		if deg == 0 {
			*h = append(*h, c)
		}
	}
	heap.Init(h)
	var res []byte
	for h.Len() > 0 {
		c := heap.Pop(h).(byte)
		res = append(res, c)
		neighbors := adj[c]
		sort.Slice(neighbors, func(i, j int) bool { return neighbors[i] < neighbors[j] })
		for _, nxt := range neighbors {
			indeg[nxt]--
			if indeg[nxt] == 0 {
				heap.Push(h, nxt)
			}
		}
	}
	if len(res) != len(indeg) {
		return ""
	}
	return string(res)
}
