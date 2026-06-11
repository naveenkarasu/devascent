class Solution {
    public long unique_paths(long m, long n) {
        long[] dp = new long[(int)n];
        java.util.Arrays.fill(dp, 1L);
        for (int i = 1; i < m; i++) {
            for (int j = 1; j < n; j++) {
                dp[j] += dp[j - 1];
            }
        }
        return dp[(int)(n - 1)];
    }
}
