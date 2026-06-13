public class Solution {
    public bool is_balanced(TreeNode root) {
        return Height(root) >= 0;
    }
    private long Height(TreeNode node) {
        if (node == null) return 0;
        long l = Height(node.left);
        if (l < 0) return -1;
        long r = Height(node.right);
        if (r < 0) return -1;
        if (Math.Abs(l - r) > 1) return -1;
        return 1 + Math.Max(l, r);
    }
}
