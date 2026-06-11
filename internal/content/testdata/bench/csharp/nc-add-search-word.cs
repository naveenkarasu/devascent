public class Solution {
    private Dictionary<char, object> root = new Dictionary<char, object>();

    private void add_word_impl(string word) {
        var node = root;
        foreach (char c in word) {
            if (!node.ContainsKey(c)) node[c] = new Dictionary<char, object>();
            node = (Dictionary<char, object>)node[c];
        }
        node['$'] = true;
    }

    private bool search_impl(Dictionary<char, object> node, string word, int i) {
        if (i == word.Length) return node.ContainsKey('$');
        char c = word[i];
        if (c == '.') {
            foreach (var kv in node) {
                if (kv.Key != '$' && search_impl((Dictionary<char, object>)kv.Value, word, i + 1))
                    return true;
            }
            return false;
        }
        if (!node.ContainsKey(c)) return false;
        return search_impl((Dictionary<char, object>)node[c], word, i + 1);
    }

    public object[] word_dictionary_ops(object[][] operations) {
        var wd_root = new Dictionary<char, object>();
        var out_list = new List<object>();
        foreach (var op in operations) {
            string cmd = (string)op[0];
            string arg = (string)op[1];
            if (cmd == "addWord") {
                var node = wd_root;
                foreach (char c in arg) {
                    if (!node.ContainsKey(c)) node[c] = new Dictionary<char, object>();
                    node = (Dictionary<char, object>)node[c];
                }
                node['$'] = true;
                out_list.Add(null);
            } else {
                out_list.Add(search_impl(wd_root, arg, 0));
            }
        }
        return out_list.ToArray();
    }
}
