class Solution {
    public long[] right_side_view(TreeNode root) {
        if (root == null) return new long[0];
        java.util.Queue<TreeNode> q = new java.util.LinkedList<>();
        java.util.List<Long> res = new java.util.ArrayList<>();
        q.offer(root);
        while (!q.isEmpty()) {
            int n = q.size();
            for (int i = 0; i < n; i++) {
                TreeNode node = q.poll();
                if (i == n - 1) res.add((long)node.val);
                if (node.left != null) q.offer(node.left);
                if (node.right != null) q.offer(node.right);
            }
        }
        long[] arr = new long[res.size()];
        for (int i = 0; i < res.size(); i++) arr[i] = res.get(i);
        return arr;
    }
}
