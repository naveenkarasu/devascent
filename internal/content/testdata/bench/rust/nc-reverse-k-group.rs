fn reverse_k_group(head: Option<Box<ListNode>>, k: i64) -> Option<Box<ListNode>> {
    // Extract values into a Vec.
    let mut vals: Vec<i32> = Vec::new();
    let mut cur = &head;
    while let Some(n) = cur {
        vals.push(n.val);
        cur = &n.next;
    }
    // Reverse each full group of k.
    let kk = k as usize;
    if kk > 0 {
        let len = vals.len();
        let mut i = 0;
        while i + kk <= len {
            vals[i..i + kk].reverse();
            i += kk;
        }
    }
    // Rebuild the list.
    let mut node: Option<Box<ListNode>> = None;
    for &v in vals.iter().rev() {
        let mut b = Box::new(ListNode::new(v));
        b.next = node;
        node = Some(b);
    }
    node
}
