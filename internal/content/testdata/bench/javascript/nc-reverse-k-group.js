function reverse_k_group(head, k) {
    function getKth(cur, k) {
        while (cur !== null && k > 0) {
            cur = cur.next;
            k--;
        }
        return cur;
    }
    const dummy = new ListNode(0, head);
    let groupPrev = dummy;
    while (true) {
        const kth = getKth(groupPrev, k);
        if (kth === null) break;
        const groupNext = kth.next;
        let prev = groupNext, cur = groupPrev.next;
        while (cur !== groupNext) {
            const nxt = cur.next;
            cur.next = prev;
            prev = cur;
            cur = nxt;
        }
        const tmp = groupPrev.next;
        groupPrev.next = kth;
        groupPrev = tmp;
    }
    return dummy.next;
}
