function reverse_k_group(head: any, k: number): any {
    function get_kth(cur: any, k: number): any {
        while (cur !== null && k > 0) {
            cur = cur.next;
            k--;
        }
        return cur;
    }

    const dummy: any = {val: 0, next: head};
    let groupPrev = dummy;

    while (true) {
        const kth = get_kth(groupPrev, k);
        if (kth === null) break;
        const groupNext = kth.next;
        let prev = groupNext;
        let cur = groupPrev.next;
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
