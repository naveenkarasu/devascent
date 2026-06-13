func add_two_numbers(l1 *ListNode, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	carry := 0
	for l1 != nil || l2 != nil || carry != 0 {
		a, b := 0, 0
		if l1 != nil {
			a = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			b = l2.Val
			l2 = l2.Next
		}
		s := a + b + carry
		carry = s / 10
		cur.Next = &ListNode{Val: s % 10}
		cur = cur.Next
	}
	return dummy.Next
}
