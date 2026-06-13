import java.util.*;

class Solution {
    static class TrieNode {
        Map<Character, TrieNode> children = new HashMap<>();
        boolean isEnd = false;
    }

    public Object[] word_dictionary_ops(String[][] operations) {
        TrieNode root = new TrieNode();
        Object[] out = new Object[operations.length];
        for (int i = 0; i < operations.length; i++) {
            String op = operations[i][0];
            String arg = operations[i][1];
            if (op.equals("addWord")) {
                TrieNode node = root;
                for (char c : arg.toCharArray()) {
                    node = node.children.computeIfAbsent(c, k -> new TrieNode());
                }
                node.isEnd = true;
                out[i] = null;
            } else { // search
                out[i] = Boolean.valueOf(dfs(root, arg, 0));
            }
        }
        return out;
    }

    private boolean dfs(TrieNode node, String word, int i) {
        if (i == word.length()) return node.isEnd;
        char c = word.charAt(i);
        if (c == '.') {
            for (TrieNode child : node.children.values()) {
                if (dfs(child, word, i + 1)) return true;
            }
            return false;
        }
        TrieNode next = node.children.get(c);
        if (next == null) return false;
        return dfs(next, word, i + 1);
    }
}
