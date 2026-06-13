class Solution {
    private int rows, cols;
    private int[][] grid;
    private static final long MOD = 1_000_000_007L;

    public long count_grid_colorings(long n, long m) {
        rows = (int)n;
        cols = (int)m;
        grid = new int[rows][cols];
        return backtrack(0);
    }

    private long backtrack(int pos) {
        if (pos == rows * cols) return 1;
        int r = pos / cols;
        int c = pos % cols;
        long total = 0;
        for (int color = 1; color <= 3; color++) {
            if ((r == 0 || grid[r - 1][c] != color) &&
                (c == 0 || grid[r][c - 1] != color)) {
                grid[r][c] = color;
                total = (total + backtrack(pos + 1)) % MOD;
                grid[r][c] = 0;
            }
        }
        return total;
    }
}
