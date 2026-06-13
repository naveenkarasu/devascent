function diameter_of_binary_tree(root: any): number {
    let best = 0;
    function height(node: any): number {
        if (node === null) return 0;
        const l = height(node.left);
        const r = height(node.right);
        if (l + r > best) best = l + r;
        return 1 + Math.max(l, r);
    }
    height(root);
    return best;
}
