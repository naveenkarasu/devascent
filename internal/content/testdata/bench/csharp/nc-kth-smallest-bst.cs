public class Solution {
    public long kth_smallest(TreeNode root, long k) {
        var stack = new Stack<TreeNode>();
        TreeNode cur = root;
        while (stack.Count > 0 || cur != null) {
            while (cur != null) {
                stack.Push(cur);
                cur = cur.left;
            }
            cur = stack.Pop();
            k--;
            if (k == 0) return cur.val;
            cur = cur.right;
        }
        return -1;
    }
}
