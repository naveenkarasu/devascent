class Solution {
    private long best;

    public long max_path_sum(TreeNode root) {
        best = Long.MIN_VALUE;
        gain(root);
        return best;
    }

    private long gain(TreeNode node) {
        if (node == null) return 0;
        long l = Math.max(gain(node.left), 0);
        long r = Math.max(gain(node.right), 0);
        best = Math.max(best, node.val + l + r);
        return node.val + Math.max(l, r);
    }
}
