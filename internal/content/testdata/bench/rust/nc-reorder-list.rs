fn reorder_list(head: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    // Collect values
    let mut vals: Vec<i32> = Vec::new();
    let mut cur = head;
    while let Some(node) = cur {
        vals.push(node.val);
        cur = node.next;
    }
    // Build reordered value sequence: L0, Ln, L1, Ln-1, ...
    let n = vals.len();
    let mut order: Vec<i32> = Vec::with_capacity(n);
    let mut i = 0usize;
    let mut j = if n == 0 { 0 } else { n - 1 };
    let mut take_front = true;
    while i <= j && order.len() < n {
        if take_front {
            order.push(vals[i]);
            if i == j {
                break;
            }
            i += 1;
        } else {
            order.push(vals[j]);
            if i == j {
                break;
            }
            j -= 1;
        }
        take_front = !take_front;
    }
    // Rebuild a fresh list from order
    let mut new_head: Option<Box<ListNode>> = None;
    for &v in order.iter().rev() {
        let mut node = Box::new(ListNode::new(v));
        node.next = new_head;
        new_head = Some(node);
    }
    new_head
}
