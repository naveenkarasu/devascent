type DetectSquares struct {
	cnt map[[2]int]int
}

func newDetectSquares() *DetectSquares {
	return &DetectSquares{cnt: make(map[[2]int]int)}
}

func (d *DetectSquares) add(x, y int) {
	d.cnt[[2]int{x, y}]++
}

func (d *DetectSquares) count(px, py int) int {
	total := 0
	for pt, c := range d.cnt {
		x, y := pt[0], pt[1]
		if abs2(x-px) == abs2(y-py) && x != px && y != py {
			total += c * d.cnt[[2]int{px, y}] * d.cnt[[2]int{x, py}]
		}
	}
	return total
}

func abs2(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func detect_squares_ops(operations [][]interface{}) []interface{} {
	ds := newDetectSquares()
	out := make([]interface{}, len(operations))
	for i, op := range operations {
		switch op[0].(string) {
		case "add":
			pt := op[1].([]interface{})
			ds.add(pt[0].(int), pt[1].(int))
			out[i] = nil
		case "count":
			pt := op[1].([]interface{})
			out[i] = ds.count(pt[0].(int), pt[1].(int))
		}
	}
	return out
}
