class Solution {
    public long[][] level_order(TreeNode root) {
        if (root == null) return new long[0][];
        java.util.Queue<TreeNode> q = new java.util.LinkedList<>();
        java.util.List<long[]> res = new java.util.ArrayList<>();
        q.offer(root);
        while (!q.isEmpty()) {
            int size = q.size();
            long[] level = new long[size];
            for (int i = 0; i < size; i++) {
                TreeNode node = q.poll();
                level[i] = node.val;
                if (node.left != null) q.offer(node.left);
                if (node.right != null) q.offer(node.right);
            }
            res.add(level);
        }
        return res.toArray(new long[0][]);
    }
}
