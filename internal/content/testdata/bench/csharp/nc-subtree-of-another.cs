public class Solution {
    public bool is_subtree(TreeNode root, TreeNode sub) {
        if (sub == null) return true;
        return Dfs(root, sub);
    }
    private bool Dfs(TreeNode node, TreeNode sub) {
        if (node == null) return false;
        if (Same(node, sub)) return true;
        return Dfs(node.left, sub) || Dfs(node.right, sub);
    }
    private bool Same(TreeNode a, TreeNode b) {
        if (a == null && b == null) return true;
        if (a == null || b == null || a.val != b.val) return false;
        return Same(a.left, b.left) && Same(a.right, b.right);
    }
}
