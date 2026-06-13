using System.Collections.Generic;

public class Solution {
    public long longest_increasing_path(long[][] matrix) {
        if (matrix == null || matrix.Length == 0 || matrix[0].Length == 0) return 0;
        int m = matrix.Length, n = matrix[0].Length;
        var memo = new Dictionary<(int, int), long>();

        long Dfs(int r, int c) {
            if (memo.TryGetValue((r, c), out long cached)) return cached;
            long best = 1;
            int[] dr = { -1, 1, 0, 0 };
            int[] dc = { 0, 0, -1, 1 };
            for (int d = 0; d < 4; d++) {
                int nr = r + dr[d], nc = c + dc[d];
                if (nr >= 0 && nr < m && nc >= 0 && nc < n && matrix[nr][nc] > matrix[r][c]) {
                    best = Math.Max(best, 1 + Dfs(nr, nc));
                }
            }
            memo[(r, c)] = best;
            return best;
        }

        long res = 0;
        for (int r = 0; r < m; r++)
            for (int c = 0; c < n; c++)
                res = Math.Max(res, Dfs(r, c));
        return res;
    }
}
