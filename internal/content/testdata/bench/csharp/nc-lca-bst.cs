public class Solution {
    public long lca_bst(TreeNode root, long p, long q) {
        TreeNode node = root;
        while (node != null) {
            if (p < node.val && q < node.val) {
                node = node.left;
            } else if (p > node.val && q > node.val) {
                node = node.right;
            } else {
                return node.val;
            }
        }
        return -1;
    }
}
