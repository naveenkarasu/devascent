type MinStack struct {
	stack []int
	mins  []int
}

func newMinStack() *MinStack {
	return &MinStack{}
}

func (s *MinStack) push(val int) {
	s.stack = append(s.stack, val)
	if len(s.mins) == 0 || val < s.mins[len(s.mins)-1] {
		s.mins = append(s.mins, val)
	} else {
		s.mins = append(s.mins, s.mins[len(s.mins)-1])
	}
}

func (s *MinStack) pop() {
	s.stack = s.stack[:len(s.stack)-1]
	s.mins = s.mins[:len(s.mins)-1]
}

func (s *MinStack) top() int {
	return s.stack[len(s.stack)-1]
}

func (s *MinStack) getMin() int {
	return s.mins[len(s.mins)-1]
}

func min_stack_ops(operations [][]interface{}) []interface{} {
	st := newMinStack()
	out := make([]interface{}, len(operations))
	for i, op := range operations {
		switch op[0].(string) {
		case "push":
			st.push(op[1].(int))
			out[i] = nil
		case "pop":
			st.pop()
			out[i] = nil
		case "top":
			out[i] = st.top()
		default: // getMin
			out[i] = st.getMin()
		}
	}
	return out
}
