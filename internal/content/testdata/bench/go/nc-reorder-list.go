func reorder_list(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	// Find middle
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	// Reverse second half
	second := slow.Next
	slow.Next = nil
	var prev *ListNode
	for second != nil {
		nxt := second.Next
		second.Next = prev
		prev = second
		second = nxt
	}
	// Merge
	first := head
	for prev != nil {
		n1 := first.Next
		n2 := prev.Next
		first.Next = prev
		prev.Next = n1
		first = n1
		prev = n2
	}
	return head
}
