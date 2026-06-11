function has_cycle(values, pos) {
    if (!values || values.length === 0) return false;
    // Build linked list internally
    const nodes = values.map(v => ({ val: v, next: null }));
    for (let i = 0; i < nodes.length - 1; i++) {
        nodes[i].next = nodes[i + 1];
    }
    if (pos >= 0) {
        nodes[nodes.length - 1].next = nodes[pos];
    }
    // Floyd's cycle detection
    let slow = nodes[0], fast = nodes[0];
    while (fast !== null && fast.next !== null) {
        slow = slow.next;
        fast = fast.next.next;
        if (slow === fast) return true;
    }
    return false;
}
