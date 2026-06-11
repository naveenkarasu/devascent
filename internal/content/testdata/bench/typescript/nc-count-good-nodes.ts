function good_nodes(root: any): number {
    function dfs(node: any, mx: number): number {
        if (node === null) return 0;
        const count = node.val >= mx ? 1 : 0;
        const newMx = Math.max(mx, node.val);
        return count + dfs(node.left, newMx) + dfs(node.right, newMx);
    }
    return dfs(root, -Infinity);
}
