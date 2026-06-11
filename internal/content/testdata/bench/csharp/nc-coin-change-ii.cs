public class Solution {
    public long change(long amount, long[] coins) {
        long[] dp = new long[amount + 1];
        dp[0] = 1;
        foreach (var coin in coins) {
            for (long x = coin; x <= amount; x++) {
                dp[x] += dp[x - coin];
            }
        }
        return dp[amount];
    }
}
