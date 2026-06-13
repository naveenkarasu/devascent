function remove_nth(head, n) {
    const dummy = new ListNode(0);
    dummy.next = head;
    let fast = dummy, slow = dummy;
    for (let i = 0; i < n; i++) {
        fast = fast.next;
    }
    while (fast.next !== null) {
        fast = fast.next;
        slow = slow.next;
    }
    slow.next = slow.next.next;
    return dummy.next;
}
