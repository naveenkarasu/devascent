public class Solution {
    private long best;
    public long diameter_of_binary_tree(TreeNode root) {
        best = 0;
        Height(root);
        return best;
    }
    private long Height(TreeNode node) {
        if (node == null) return 0;
        long l = Height(node.left);
        long r = Height(node.right);
        if (l + r > best) best = l + r;
        return 1 + Math.Max(l, r);
    }
}
