fn add_two_numbers(l1: Option<Box<ListNode>>, l2: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    let mut a = &l1;
    let mut b = &l2;
    let mut carry = 0i64;
    let mut digits: Vec<i32> = Vec::new();
    loop {
        let av = match a {
            Some(n) => n.val as i64,
            None => 0,
        };
        let bv = match b {
            Some(n) => n.val as i64,
            None => 0,
        };
        if a.is_none() && b.is_none() && carry == 0 {
            break;
        }
        let s = av + bv + carry;
        carry = s / 10;
        digits.push((s % 10) as i32);
        if let Some(n) = a {
            a = &n.next;
        }
        if let Some(n) = b {
            b = &n.next;
        }
    }
    // Rebuild in order (digits are already in result order, least-significant first).
    let mut node: Option<Box<ListNode>> = None;
    for &d in digits.iter().rev() {
        let mut bx = Box::new(ListNode::new(d));
        bx.next = node;
        node = Some(bx);
    }
    node
}
