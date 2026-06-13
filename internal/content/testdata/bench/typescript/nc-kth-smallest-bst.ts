function kth_smallest(root: any, k: number): number {
    const stack: any[] = [];
    let cur = root;
    while (stack.length > 0 || cur !== null) {
        while (cur !== null) {
            stack.push(cur);
            cur = cur.left;
        }
        cur = stack.pop();
        k -= 1;
        if (k === 0) return cur.val;
        cur = cur.right;
    }
    return -1;
}
