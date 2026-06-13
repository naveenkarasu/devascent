public class Solution {
    public string[] find_words(string[][] board, string[] words) {
        var trie = new Dictionary<char, object>();
        foreach (string w in words) {
            var node = trie;
            foreach (char c in w) {
                if (!node.ContainsKey(c)) node[c] = new Dictionary<char, object>();
                node = (Dictionary<char, object>)node[c];
            }
            node['$'] = w;
        }

        if (board == null || board.Length == 0 || board[0].Length == 0) return new string[0];
        int rows = board.Length, cols = board[0].Length;
        var found = new HashSet<string>();
        char[][] b = new char[rows][];
        for (int r = 0; r < rows; r++) {
            b[r] = new char[cols];
            for (int c = 0; c < cols; c++) b[r][c] = board[r][c][0];
        }

        void dfs(int r, int c, Dictionary<char, object> node) {
            char ch = b[r][c];
            if (!node.ContainsKey(ch)) return;
            var nxt = (Dictionary<char, object>)node[ch];
            if (nxt.ContainsKey('$')) found.Add((string)nxt['$']);
            b[r][c] = '#';
            int[][] dirs = { new[]{1,0}, new[]{-1,0}, new[]{0,1}, new[]{0,-1} };
            foreach (var d in dirs) {
                int nr = r + d[0], nc = c + d[1];
                if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && b[nr][nc] != '#')
                    dfs(nr, nc, nxt);
            }
            b[r][c] = ch;
        }

        for (int r = 0; r < rows; r++)
            for (int c = 0; c < cols; c++)
                dfs(r, c, trie);

        var result = new List<string>(found);
        result.Sort();
        return result.ToArray();
    }
}
