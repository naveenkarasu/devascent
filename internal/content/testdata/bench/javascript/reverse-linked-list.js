function reverse_list(head) {
    let prev = null;
    while (head !== null) {
        const nxt = head.next;
        head.next = prev;
        prev = head;
        head = nxt;
    }
    return prev;
}
