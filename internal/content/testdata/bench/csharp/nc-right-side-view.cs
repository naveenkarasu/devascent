public class Solution {
    public long[] right_side_view(TreeNode root) {
        if (root == null) return new long[0];
        var q = new Queue<TreeNode>();
        var res = new List<long>();
        q.Enqueue(root);
        while (q.Count > 0) {
            int n = q.Count;
            for (int i = 0; i < n; i++) {
                TreeNode node = q.Dequeue();
                if (i == n - 1) res.Add(node.val);
                if (node.left != null) q.Enqueue(node.left);
                if (node.right != null) q.Enqueue(node.right);
            }
        }
        return res.ToArray();
    }
}
