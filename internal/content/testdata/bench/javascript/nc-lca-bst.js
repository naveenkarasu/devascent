function lca_bst(root, p, q) {
    let node = root;
    while (node !== null) {
        if (p < node.val && q < node.val) {
            node = node.left;
        } else if (p > node.val && q > node.val) {
            node = node.right;
        } else {
            return node.val;
        }
    }
    return -1;
}
