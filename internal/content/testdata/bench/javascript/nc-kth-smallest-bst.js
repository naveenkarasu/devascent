function kth_smallest(root, k) {
    const stack = [];
    let cur = root;
    while (stack.length > 0 || cur !== null) {
        while (cur !== null) {
            stack.push(cur);
            cur = cur.left;
        }
        cur = stack.pop();
        k--;
        if (k === 0) return cur.val;
        cur = cur.right;
    }
    return -1;
}
