fn merge_lists(a: Option<Box<ListNode>>, b: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    match (a, b) {
        (None, b) => b,
        (a, None) => a,
        (Some(mut na), Some(mut nb)) => {
            if na.val <= nb.val {
                na.next = merge_lists(na.next.take(), Some(nb));
                Some(na)
            } else {
                nb.next = merge_lists(Some(na), nb.next.take());
                Some(nb)
            }
        }
    }
}
