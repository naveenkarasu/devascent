class Solution {
    public long min_distance(String word1, String word2) {
        int m = word1.length(), n = word2.length();
        long[] dp = new long[n + 1];
        for (int j = 0; j <= n; j++) dp[j] = j;
        for (int i = 1; i <= m; i++) {
            long prev = dp[0];
            dp[0] = i;
            for (int j = 1; j <= n; j++) {
                long temp = dp[j];
                if (word1.charAt(i - 1) == word2.charAt(j - 1)) {
                    dp[j] = prev;
                } else {
                    dp[j] = 1 + Math.min(prev, Math.min(dp[j], dp[j - 1]));
                }
                prev = temp;
            }
        }
        return dp[n];
    }
}
