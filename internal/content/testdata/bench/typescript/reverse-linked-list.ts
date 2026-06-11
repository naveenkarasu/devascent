function reverse_list(head: any): any {
    let prev: any = null;
    while (head !== null) {
        const nxt = head.next;
        head.next = prev;
        prev = head;
        head = nxt;
    }
    return prev;
}
