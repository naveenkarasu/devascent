func copy_list(nodes [][]interface{}) [][]interface{} {
	n := len(nodes)
	if n == 0 {
		return [][]interface{}{}
	}
	val := make([]int, n)
	rnd := make([]int, n) // index, or -1 for nil
	for i := range nodes {
		val[i] = nodes[i][0].(int)
		if nodes[i][1] == nil {
			rnd[i] = -1
		} else {
			rnd[i] = nodes[i][1].(int)
		}
	}
	cval := make([]int, n)
	crnd := make([]int, n)
	copy(cval, val)
	copy(crnd, rnd)
	out := make([][]interface{}, n)
	for i := 0; i < n; i++ {
		var r interface{}
		if crnd[i] >= 0 {
			r = crnd[i]
		} else {
			r = nil
		}
		out[i] = []interface{}{cval[i], r}
	}
	return out
}
