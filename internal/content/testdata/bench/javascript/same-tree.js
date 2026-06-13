function is_same_tree(p, q) {
    if (p === null && q === null) return true;
    if (p === null || q === null || p.val !== q.val) return false;
    return is_same_tree(p.left, q.left) && is_same_tree(p.right, q.right);
}
