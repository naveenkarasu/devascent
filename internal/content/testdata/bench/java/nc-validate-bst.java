class Solution {
    public boolean is_valid_bst(TreeNode root) {
        return valid(root, Long.MIN_VALUE, Long.MAX_VALUE);
    }
    private boolean valid(TreeNode node, long lo, long hi) {
        if (node == null) return true;
        if (!(lo < node.val && node.val < hi)) return false;
        return valid(node.left, lo, node.val) && valid(node.right, node.val, hi);
    }
}
