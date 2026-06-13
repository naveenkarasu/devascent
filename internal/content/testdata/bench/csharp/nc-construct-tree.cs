public class Solution {
    public TreeNode build_tree(long[] preorder, long[] inorder) {
        var idx = new Dictionary<long, int>();
        for (int i = 0; i < inorder.Length; i++) idx[inorder[i]] = i;
        int pre = 0;
        TreeNode Build(int lo, int hi) {
            if (lo > hi) return null;
            long val = preorder[pre++];
            var node = new TreeNode((int)val);
            int mid = idx[val];
            node.left = Build(lo, mid - 1);
            node.right = Build(mid + 1, hi);
            return node;
        }
        return Build(0, inorder.Length - 1);
    }
}
