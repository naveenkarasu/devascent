function copy_list(arr: any[][]): any[][] {
    if (!arr || arr.length === 0) return [];

    // Build original nodes
    const nodes: any[] = arr.map(([v]) => ({val: v, next: null, random: null}));
    for (let i = 0; i < nodes.length; i++) {
        nodes[i].next = i + 1 < nodes.length ? nodes[i + 1] : null;
        const r = arr[i][1];
        nodes[i].random = r !== null && r !== undefined ? nodes[r] : null;
    }

    // Deep copy
    const mp = new Map<any, any>();
    let cur = nodes[0];
    while (cur !== null) {
        mp.set(cur, {val: cur.val, next: null, random: null});
        cur = cur.next;
    }
    cur = nodes[0];
    while (cur !== null) {
        const copy = mp.get(cur);
        copy.next = cur.next ? mp.get(cur.next) : null;
        copy.random = cur.random ? mp.get(cur.random) : null;
        cur = cur.next;
    }

    // Build output: index map
    const copyHead = mp.get(nodes[0]);
    const orderMap = new Map<any, number>();
    let c = copyHead;
    let idx = 0;
    while (c !== null) {
        orderMap.set(c, idx);
        c = c.next;
        idx++;
    }

    const out: any[][] = [];
    c = copyHead;
    while (c !== null) {
        const randomIdx = c.random !== null ? orderMap.get(c.random)! : null;
        out.push([c.val, randomIdx]);
        c = c.next;
    }
    return out;
}
