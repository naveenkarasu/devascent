function build_tree(preorder, inorder) {
    const idx = new Map();
    for (let i = 0; i < inorder.length; i++) {
        idx.set(inorder[i], i);
    }
    let pre = 0;
    function build(lo, hi) {
        if (lo > hi) return null;
        const val = preorder[pre++];
        const node = { val: val, left: null, right: null };
        const mid = idx.get(val);
        node.left = build(lo, mid - 1);
        node.right = build(mid + 1, hi);
        return node;
    }
    return build(0, inorder.length - 1);
}
