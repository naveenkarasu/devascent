function add_two_numbers(l1, l2) {
    const dummy = new ListNode(0);
    let cur = dummy;
    let carry = 0;
    while (l1 !== null || l2 !== null || carry) {
        const a = l1 !== null ? l1.val : 0;
        const b = l2 !== null ? l2.val : 0;
        const s = a + b + carry;
        carry = Math.floor(s / 10);
        cur.next = new ListNode(s % 10);
        cur = cur.next;
        if (l1 !== null) l1 = l1.next;
        if (l2 !== null) l2 = l2.next;
    }
    return dummy.next;
}
