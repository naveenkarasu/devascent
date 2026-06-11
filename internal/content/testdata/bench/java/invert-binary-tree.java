class Solution {
    public TreeNode invert_tree(TreeNode root) {
        if (root == null) return null;
        TreeNode tmp = root.left;
        root.left = invert_tree(root.right);
        root.right = invert_tree(tmp);
        return root;
    }
}
