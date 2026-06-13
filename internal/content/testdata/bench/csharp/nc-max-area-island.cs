public class Solution {
    public long max_area_of_island(long[][] grid) {
        if (grid == null || grid.Length == 0) return 0;
        int m = grid.Length, n = grid[0].Length;
        bool[][] seen = new bool[m][];
        for (int i = 0; i < m; i++) seen[i] = new bool[n];
        long best = 0;
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (grid[i][j] == 1 && !seen[i][j]) {
                    best = Math.Max(best, Dfs(grid, seen, i, j, m, n));
                }
            }
        }
        return best;
    }

    private long Dfs(long[][] grid, bool[][] seen, int i, int j, int m, int n) {
        if (i < 0 || j < 0 || i >= m || j >= n) return 0;
        if (grid[i][j] == 0 || seen[i][j]) return 0;
        seen[i][j] = true;
        return 1 + Dfs(grid, seen, i+1, j, m, n)
                 + Dfs(grid, seen, i-1, j, m, n)
                 + Dfs(grid, seen, i, j+1, m, n)
                 + Dfs(grid, seen, i, j-1, m, n);
    }
}
