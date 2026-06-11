class Solution {
    public long num_islands(long[][] grid) {
        if (grid == null || grid.length == 0) return 0;
        int m = grid.length, n = grid[0].length;
        boolean[][] seen = new boolean[m][n];
        long count = 0;
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (grid[i][j] == 1 && !seen[i][j]) {
                    count++;
                    dfs(grid, seen, i, j, m, n);
                }
            }
        }
        return count;
    }
    private void dfs(long[][] grid, boolean[][] seen, int i, int j, int m, int n) {
        if (i < 0 || j < 0 || i >= m || j >= n) return;
        if (grid[i][j] == 0 || seen[i][j]) return;
        seen[i][j] = true;
        dfs(grid, seen, i + 1, j, m, n);
        dfs(grid, seen, i - 1, j, m, n);
        dfs(grid, seen, i, j + 1, m, n);
        dfs(grid, seen, i, j - 1, m, n);
    }
}
