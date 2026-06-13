fn reverse_list(head: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    let mut prev: Option<Box<ListNode>> = None;
    let mut cur = head;
    while let Some(mut node) = cur {
        cur = node.next.take();
        node.next = prev;
        prev = Some(node);
    }
    prev
}
