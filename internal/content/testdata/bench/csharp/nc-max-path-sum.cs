public class Solution {
    public long max_path_sum(TreeNode root) {
        long best = long.MinValue;
        long Gain(TreeNode node) {
            if (node == null) return 0;
            long l = Math.Max(Gain(node.left), 0);
            long r = Math.Max(Gain(node.right), 0);
            best = Math.Max(best, node.val + l + r);
            return node.val + Math.Max(l, r);
        }
        Gain(root);
        return best;
    }
}
