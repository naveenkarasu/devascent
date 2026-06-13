function is_subtree(root, sub) {
    function same(a, b) {
        if (a === null && b === null) return true;
        if (a === null || b === null || a.val !== b.val) return false;
        return same(a.left, b.left) && same(a.right, b.right);
    }
    if (sub === null) return true;
    function dfs(node) {
        if (node === null) return false;
        if (same(node, sub)) return true;
        return dfs(node.left) || dfs(node.right);
    }
    return dfs(root);
}
