function lca_bst(root: any, p: number, q: number): number {
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
