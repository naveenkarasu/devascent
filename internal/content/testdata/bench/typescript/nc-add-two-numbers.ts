function add_two_numbers(l1: any, l2: any): any {
    const dummy: any = {val: 0, next: null};
    let cur = dummy;
    let carry = 0;
    while (l1 !== null || l2 !== null || carry) {
        const a = l1 !== null ? l1.val : 0;
        const b = l2 !== null ? l2.val : 0;
        const s = a + b + carry;
        carry = Math.floor(s / 10);
        cur.next = {val: s % 10, next: null};
        cur = cur.next;
        if (l1 !== null) l1 = l1.next;
        if (l2 !== null) l2 = l2.next;
    }
    return dummy.next;
}
