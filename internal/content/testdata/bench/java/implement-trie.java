import java.util.*;

class Solution {
    static class TrieNode {
        Map<Character, TrieNode> children = new HashMap<>();
        boolean isEnd = false;
    }

    public Object[] trie_ops(String[][] operations) {
        TrieNode root = new TrieNode();
        Object[] out = new Object[operations.length];
        for (int i = 0; i < operations.length; i++) {
            String op = operations[i][0];
            String arg = operations[i][1];
            if (op.equals("insert")) {
                TrieNode node = root;
                for (char c : arg.toCharArray()) {
                    node = node.children.computeIfAbsent(c, k -> new TrieNode());
                }
                node.isEnd = true;
                out[i] = null;
            } else if (op.equals("search")) {
                TrieNode node = find(root, arg);
                out[i] = Boolean.valueOf(node != null && node.isEnd);
            } else { // startsWith
                TrieNode node = find(root, arg);
                out[i] = Boolean.valueOf(node != null);
            }
        }
        return out;
    }

    private TrieNode find(TrieNode root, String word) {
        TrieNode node = root;
        for (char c : word.toCharArray()) {
            node = node.children.get(c);
            if (node == null) return null;
        }
        return node;
    }
}
