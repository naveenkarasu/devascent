function invert_tree(root) {
    if (root === null) return null;
    const tmp = root.left;
    root.left = invert_tree(root.right);
    root.right = invert_tree(tmp);
    return root;
}
