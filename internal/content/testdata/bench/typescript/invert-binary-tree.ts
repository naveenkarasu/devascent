function invert_tree(root: any): any {
    if (root === null) return null;
    const left = invert_tree(root.right);
    const right = invert_tree(root.left);
    root.left = left;
    root.right = right;
    return root;
}
