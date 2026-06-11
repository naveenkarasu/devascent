public class Solution {
    public long[] counting_bits(long n) {
        long[] dp = new long[n + 1];
        for (long i = 1; i <= n; i++) {
            dp[i] = dp[i >> 1] + (i & 1);
        }
        return dp;
    }
}
