class Solution {
    public long max_hourglass_sum(long[][] grid) {
        long best = Long.MIN_VALUE;
        for (int i = 0; i < 4; i++) {
            for (int j = 0; j < 4; j++) {
                long s = grid[i][j] + grid[i][j+1] + grid[i][j+2]
                       + grid[i+1][j+1]
                       + grid[i+2][j] + grid[i+2][j+1] + grid[i+2][j+2];
                if (s > best) best = s;
            }
        }
        return best;
    }
}
