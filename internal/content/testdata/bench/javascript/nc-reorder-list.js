function reorder_list(head) {
    if (!head || !head.next) return head;
    // Find middle
    let slow = head, fast = head;
    while (fast.next !== null && fast.next.next !== null) {
        slow = slow.next;
        fast = fast.next.next;
    }
    // Reverse second half
    let second = slow.next;
    slow.next = null;
    let prev = null;
    while (second !== null) {
        const nxt = second.next;
        second.next = prev;
        prev = second;
        second = nxt;
    }
    // Merge two halves
    let first = head;
    while (prev !== null) {
        const n1 = first.next;
        const n2 = prev.next;
        first.next = prev;
        prev.next = n1;
        first = n1;
        prev = n2;
    }
    return head;
}
