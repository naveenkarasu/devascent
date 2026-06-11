public class Solution {
    public long[][] level_order(TreeNode root) {
        if (root == null) return new long[0][];
        var q = new Queue<TreeNode>();
        var res = new List<long[]>();
        q.Enqueue(root);
        while (q.Count > 0) {
            int size = q.Count;
            long[] level = new long[size];
            for (int i = 0; i < size; i++) {
                TreeNode node = q.Dequeue();
                level[i] = node.val;
                if (node.left != null) q.Enqueue(node.left);
                if (node.right != null) q.Enqueue(node.right);
            }
            res.Add(level);
        }
        return res.ToArray();
    }
}
