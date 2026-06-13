import java.util.Arrays;

class Solution {
    public long coin_change(long[] coins, long amount) {
        int amt = (int) amount;
        long INF = amt + 1;
        long[] dp = new long[amt + 1];
        Arrays.fill(dp, INF);
        dp[0] = 0;
        for (int a = 1; a <= amt; a++) {
            for (long c : coins) {
                if (c <= a) {
                    dp[a] = Math.min(dp[a], dp[(int)(a - c)] + 1);
                }
            }
        }
        return dp[amt] != INF ? dp[amt] : -1;
    }
}
