public class Solution {
    public long good_nodes(TreeNode root) {
        return Dfs(root, long.MinValue);
    }

    private long Dfs(TreeNode node, long mx) {
        if (node == null) return 0;
        long count = (node.val >= mx) ? 1 : 0;
        long newMx = Math.Max(mx, node.val);
        return count + Dfs(node.left, newMx) + Dfs(node.right, newMx);
    }
}
