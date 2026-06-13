function build_tree(preorder: number[], inorder: number[]): any {
    const idx: Record<number, number> = {};
    for (let i = 0; i < inorder.length; i++) {
        idx[inorder[i]] = i;
    }
    let pre = 0;
    function build(lo: number, hi: number): any {
        if (lo > hi) return null;
        const val = preorder[pre++];
        const node = { val, left: null, right: null };
        const mid = idx[val];
        node.left = build(lo, mid - 1);
        node.right = build(mid + 1, hi);
        return node;
    }
    return build(0, inorder.length - 1);
}
