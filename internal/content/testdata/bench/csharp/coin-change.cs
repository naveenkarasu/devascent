public class Solution {
    public long coin_change(long[] coins, long amount) {
        long INF = amount + 1;
        var dp = new long[amount + 1];
        for (int i = 1; i <= amount; i++) dp[i] = INF;
        for (long a = 1; a <= amount; a++) {
            foreach (long c in coins) {
                if (c <= a) dp[a] = Math.Min(dp[a], dp[a - c] + 1);
            }
        }
        return dp[amount] != INF ? dp[amount] : -1;
    }
}
