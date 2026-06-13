class Solution {
    public boolean is_balanced(TreeNode root) {
        return height(root) >= 0;
    }
    private long height(TreeNode node) {
        if (node == null) return 0;
        long l = height(node.left);
        if (l < 0) return -1;
        long r = height(node.right);
        if (r < 0) return -1;
        if (Math.abs(l - r) > 1) return -1;
        return 1 + Math.max(l, r);
    }
}
