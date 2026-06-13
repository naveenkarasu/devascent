import java.util.*;

class Solution {
    public String[] find_words(String[][] board, String[] words) {
        // Build trie
        Map<String, Object> trie = new HashMap<>();
        for (String w : words) {
            Map<String, Object> node = trie;
            for (char c : w.toCharArray()) {
                String key = String.valueOf(c);
                node = (Map<String, Object>) node.computeIfAbsent(key, k -> new HashMap<>());
            }
            node.put("$", w);
        }
        if (board == null || board.length == 0) return new String[0];
        int rows = board.length, cols = board[0].length;
        Set<String> found = new HashSet<>();
        boolean[][] visited = new boolean[rows][cols];

        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                dfs(board, r, c, rows, cols, trie, found, visited);
            }
        }
        String[] res = found.toArray(new String[0]);
        Arrays.sort(res);
        return res;
    }

    @SuppressWarnings("unchecked")
    private void dfs(String[][] board, int r, int c, int rows, int cols,
                     Map<String, Object> node, Set<String> found, boolean[][] visited) {
        if (r < 0 || r >= rows || c < 0 || c >= cols || visited[r][c]) return;
        String key = board[r][c];
        if (!node.containsKey(key)) return;
        Map<String, Object> nxt = (Map<String, Object>) node.get(key);
        if (nxt.containsKey("$")) found.add((String) nxt.get("$"));
        visited[r][c] = true;
        int[][] dirs = {{1,0},{-1,0},{0,1},{0,-1}};
        for (int[] d : dirs) {
            dfs(board, r + d[0], c + d[1], rows, cols, nxt, found, visited);
        }
        visited[r][c] = false;
    }
}
