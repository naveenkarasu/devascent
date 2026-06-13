function is_balanced(root) {
    function height(node) {
        if (node === null) return 0;
        const l = height(node.left);
        if (l < 0) return -1;
        const r = height(node.right);
        if (r < 0) return -1;
        if (Math.abs(l - r) > 1) return -1;
        return 1 + Math.max(l, r);
    }
    return height(root) >= 0;
}
