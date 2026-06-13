class Solution {
    public boolean is_subtree(TreeNode root, TreeNode sub) {
        if (sub == null) return true;
        return dfs(root, sub);
    }
    private boolean dfs(TreeNode node, TreeNode sub) {
        if (node == null) return false;
        if (same(node, sub)) return true;
        return dfs(node.left, sub) || dfs(node.right, sub);
    }
    private boolean same(TreeNode a, TreeNode b) {
        if (a == null && b == null) return true;
        if (a == null || b == null || a.val != b.val) return false;
        return same(a.left, b.left) && same(a.right, b.right);
    }
}
