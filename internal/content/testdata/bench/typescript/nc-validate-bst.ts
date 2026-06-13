function is_valid_bst(root: any): boolean {
    function valid(node: any, lo: number, hi: number): boolean {
        if (node === null) return true;
        if (!(lo < node.val && node.val < hi)) return false;
        return valid(node.left, lo, node.val) && valid(node.right, node.val, hi);
    }
    return valid(root, -Infinity, Infinity);
}
