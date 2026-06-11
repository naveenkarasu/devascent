class Solution {
    public TreeNode build_tree(long[] preorder, long[] inorder) {
        java.util.Map<Long, Integer> idx = new java.util.HashMap<>();
        for (int i = 0; i < inorder.length; i++) idx.put(inorder[i], i);
        int[] pre = {0};
        return build(preorder, idx, pre, 0, inorder.length - 1);
    }

    private TreeNode build(long[] preorder, java.util.Map<Long, Integer> idx, int[] pre, int lo, int hi) {
        if (lo > hi) return null;
        long val = preorder[pre[0]++];
        TreeNode node = new TreeNode((int) val);
        int mid = idx.get(val);
        node.left = build(preorder, idx, pre, lo, mid - 1);
        node.right = build(preorder, idx, pre, mid + 1, hi);
        return node;
    }
}
