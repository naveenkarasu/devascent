func reverse_k_group(head *ListNode, k int) *ListNode {
	getKth := func(cur *ListNode, k int) *ListNode {
		for cur != nil && k > 0 {
			cur = cur.Next
			k--
		}
		return cur
	}
	dummy := &ListNode{Next: head}
	groupPrev := dummy
	for {
		kth := getKth(groupPrev, k)
		if kth == nil {
			break
		}
		groupNext := kth.Next
		prev, cur := groupNext, groupPrev.Next
		for cur != groupNext {
			nxt := cur.Next
			cur.Next = prev
			prev = cur
			cur = nxt
		}
		tmp := groupPrev.Next
		groupPrev.Next = kth
		groupPrev = tmp
	}
	return dummy.Next
}
