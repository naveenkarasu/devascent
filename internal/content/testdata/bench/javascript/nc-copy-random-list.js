function copy_list(arr) {
    if (!arr || arr.length === 0) return [];
    // Build original nodes
    const nodes = arr.map(([v]) => ({ val: v, next: null, random: null }));
    for (let i = 0; i < nodes.length; i++) {
        nodes[i].next = i + 1 < nodes.length ? nodes[i + 1] : null;
        nodes[i].random = arr[i][1] !== null ? nodes[arr[i][1]] : null;
    }
    // Deep copy
    const mp = new Map();
    for (const node of nodes) {
        mp.set(node, { val: node.val, next: null, random: null });
    }
    for (const node of nodes) {
        const copy = mp.get(node);
        copy.next = node.next ? mp.get(node.next) : null;
        copy.random = node.random ? mp.get(node.random) : null;
    }
    const copyHead = mp.get(nodes[0]);
    // Build index map for copy nodes
    const order = new Map();
    let cur = copyHead, idx = 0;
    while (cur !== null) {
        order.set(cur, idx++);
        cur = cur.next;
    }
    // Output
    const out = [];
    cur = copyHead;
    while (cur !== null) {
        out.push([cur.val, cur.random !== null ? order.get(cur.random) : null]);
        cur = cur.next;
    }
    return out;
}
