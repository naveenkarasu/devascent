function remove_nth(head: any, n: number): any {
    const dummy: any = { val: 0, next: head };
    let fast: any = dummy, slow: any = dummy;
    for (let i = 0; i < n; i++) fast = fast.next;
    while (fast.next !== null) {
        fast = fast.next;
        slow = slow.next;
    }
    slow.next = slow.next.next;
    return dummy.next;
}
