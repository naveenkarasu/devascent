public class Solution {
    public long num_decodings(string s) {
        int n = s.Length;
        long[] dp = new long[n + 1];
        dp[0] = 1;
        dp[1] = s[0] == '0' ? 0 : 1;
        for (int i = 2; i <= n; i++) {
            int one = s[i - 1] - '0';
            int two = int.Parse(s.Substring(i - 2, 2));
            if (one != 0) dp[i] += dp[i - 1];
            if (two >= 10 && two <= 26) dp[i] += dp[i - 2];
        }
        return dp[n];
    }
}
