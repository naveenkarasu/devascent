class Solution {
    private long best = 0;
    public long diameter_of_binary_tree(TreeNode root) {
        best = 0;
        height(root);
        return best;
    }
    private long height(TreeNode node) {
        if (node == null) return 0;
        long l = height(node.left);
        long r = height(node.right);
        if (l + r > best) best = l + r;
        return 1 + Math.max(l, r);
    }
}
