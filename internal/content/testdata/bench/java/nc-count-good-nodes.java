class Solution {
    public long good_nodes(TreeNode root) {
        return dfs(root, Integer.MIN_VALUE);
    }

    private long dfs(TreeNode node, int mx) {
        if (node == null) return 0;
        long count = (node.val >= mx) ? 1 : 0;
        int newMx = Math.max(mx, node.val);
        return count + dfs(node.left, newMx) + dfs(node.right, newMx);
    }
}
