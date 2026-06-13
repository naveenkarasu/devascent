public class Solution {
    public long max_depth(TreeNode root) {
        if (root == null) return 0;
        return 1 + Math.Max(max_depth(root.left), max_depth(root.right));
    }
}
