public class Solution {
    private string Serialize(TreeNode root) {
        var out_list = new List<string>();
        void Dfs(TreeNode node) {
            if (node == null) { out_list.Add("#"); return; }
            out_list.Add(node.val.ToString());
            Dfs(node.left);
            Dfs(node.right);
        }
        Dfs(root);
        return string.Join(",", out_list);
    }

    private TreeNode Deserialize(string data) {
        var vals = new Queue<string>(data.Split(','));
        TreeNode Build() {
            string v = vals.Dequeue();
            if (v == "#") return null;
            var node = new TreeNode(int.Parse(v));
            node.left = Build();
            node.right = Build();
            return node;
        }
        return Build();
    }

    public TreeNode codec_roundtrip(TreeNode root) {
        string data = Serialize(root);
        return Deserialize(data);
    }
}
