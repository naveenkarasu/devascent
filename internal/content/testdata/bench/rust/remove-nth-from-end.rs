fn remove_nth(head: Option<Box<ListNode>>, n: i64) -> Option<Box<ListNode>> {
    let mut vals: Vec<i32> = Vec::new();
    let mut cur = &head;
    while let Some(node) = cur {
        vals.push(node.val);
        cur = &node.next;
    }
    let len = vals.len();
    let idx = len - n as usize; // index to drop from the front
    let mut out: Option<Box<ListNode>> = None;
    for (i, &v) in vals.iter().enumerate().rev() {
        if i == idx {
            continue;
        }
        out = Some(Box::new(ListNode { val: v, next: out }));
    }
    out
}
