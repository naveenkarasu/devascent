class Solution {
    public long[] counting_bits(long n) {
        long[] dp = new long[(int)n + 1];
        for (int i = 1; i <= (int)n; i++) {
            dp[i] = dp[i >> 1] + (i & 1);
        }
        return dp;
    }
}
