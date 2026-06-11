public class Solution {
    public TreeNode invert_tree(TreeNode root) {
        if (root == null) return null;
        TreeNode left = invert_tree(root.right);
        TreeNode right = invert_tree(root.left);
        root.left = left;
        root.right = right;
        return root;
    }
}
