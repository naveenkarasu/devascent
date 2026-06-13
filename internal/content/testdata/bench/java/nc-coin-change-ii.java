class Solution {
    public long change(long amount, long[] coins) {
        long[] dp = new long[(int)amount + 1];
        dp[0] = 1;
        for (long coin : coins) {
            for (long x = coin; x <= amount; x++) {
                dp[(int)x] += dp[(int)(x - coin)];
            }
        }
        return dp[(int)amount];
    }
}
