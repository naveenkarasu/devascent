public class Solution {
    private class TrieNode {
        public Dictionary<char, TrieNode> Children = new Dictionary<char, TrieNode>();
        public bool IsEnd = false;
    }

    public object[] trie_ops(object[][] operations) {
        var root = new TrieNode();
        var out_ = new object[operations.Length];
        for (int i = 0; i < operations.Length; i++) {
            string op = (string)operations[i][0];
            string arg = (string)operations[i][1];
            if (op == "insert") {
                var node = root;
                foreach (char c in arg) {
                    if (!node.Children.ContainsKey(c))
                        node.Children[c] = new TrieNode();
                    node = node.Children[c];
                }
                node.IsEnd = true;
                out_[i] = null;
            } else if (op == "search") {
                var node = Find(root, arg);
                out_[i] = node != null && node.IsEnd;
            } else { // startsWith
                out_[i] = Find(root, arg) != null;
            }
        }
        return out_;
    }

    private TrieNode Find(TrieNode root, string word) {
        var node = root;
        foreach (char c in word) {
            if (!node.Children.ContainsKey(c)) return null;
            node = node.Children[c];
        }
        return node;
    }
}
