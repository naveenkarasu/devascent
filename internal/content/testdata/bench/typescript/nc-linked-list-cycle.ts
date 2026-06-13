function has_cycle(values: number[], pos: number): boolean {
    if (values.length === 0) return false;
    const nodes: any[] = values.map(v => ({ val: v, next: null }));
    for (let i = 0; i < nodes.length - 1; i++) {
        nodes[i].next = nodes[i + 1];
    }
    if (pos >= 0 && nodes.length > 0) {
        nodes[nodes.length - 1].next = nodes[pos];
    }
    let slow: any = nodes[0], fast: any = nodes[0];
    while (fast !== null && fast.next !== null) {
        slow = slow.next;
        fast = fast.next.next;
        if (slow === fast) return true;
    }
    return false;
}
