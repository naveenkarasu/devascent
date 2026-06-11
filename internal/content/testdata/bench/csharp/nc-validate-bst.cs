public class Solution {
    public bool is_valid_bst(TreeNode root) {
        return Valid(root, long.MinValue, long.MaxValue);
    }
    private bool Valid(TreeNode node, long lo, long hi) {
        if (node == null) return true;
        if (!(lo < node.val && node.val < hi)) return false;
        return Valid(node.left, lo, node.val) && Valid(node.right, node.val, hi);
    }
}
