class Solution {
    private String serialize(TreeNode root) {
        StringBuilder sb = new StringBuilder();
        serializeHelper(root, sb);
        return sb.toString();
    }

    private void serializeHelper(TreeNode node, StringBuilder sb) {
        if (sb.length() > 0) sb.append(',');
        if (node == null) { sb.append('#'); return; }
        sb.append(node.val);
        serializeHelper(node.left, sb);
        serializeHelper(node.right, sb);
    }

    private TreeNode deserialize(String data) {
        java.util.Queue<String> q = new java.util.LinkedList<>(java.util.Arrays.asList(data.split(",")));
        return buildTree(q);
    }

    private TreeNode buildTree(java.util.Queue<String> q) {
        String v = q.poll();
        if ("#".equals(v)) return null;
        TreeNode node = new TreeNode(Integer.parseInt(v));
        node.left = buildTree(q);
        node.right = buildTree(q);
        return node;
    }

    public TreeNode codec_roundtrip(TreeNode root) {
        return deserialize(serialize(root));
    }
}
