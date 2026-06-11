public class Solution {
    public long edit_distance(string word1, string word2) {
        int m = word1.Length, n = word2.Length;
        long[] dp = new long[n + 1];
        for (int j = 0; j <= n; j++) dp[j] = j;
        for (int i = 1; i <= m; i++) {
            long prev = dp[0];
            dp[0] = i;
            for (int j = 1; j <= n; j++) {
                long temp = dp[j];
                if (word1[i - 1] == word2[j - 1]) {
                    dp[j] = prev;
                } else {
                    dp[j] = 1 + Math.Min(prev, Math.Min(dp[j - 1], dp[j]));
                }
                prev = temp;
            }
        }
        return dp[n];
    }
}
