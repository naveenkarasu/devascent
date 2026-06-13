function merge_lists(a, b) {
    const dummy = new ListNode();
    let tail = dummy;
    while (a !== null && b !== null) {
        if (a.val <= b.val) {
            tail.next = a;
            a = a.next;
        } else {
            tail.next = b;
            b = b.next;
        }
        tail = tail.next;
    }
    tail.next = (a !== null) ? a : b;
    return dummy.next;
}
