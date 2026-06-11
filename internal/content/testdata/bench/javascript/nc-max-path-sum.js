function max_path_sum(root) {
    let best = [-Infinity];
    function gain(node) {
        if (node === null) return 0;
        const l = Math.max(gain(node.left), 0);
        const r = Math.max(gain(node.right), 0);
        best[0] = Math.max(best[0], node.val + l + r);
        return node.val + Math.max(l, r);
    }
    gain(root);
    return best[0];
}
