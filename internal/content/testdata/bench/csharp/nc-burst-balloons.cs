public class Solution {
    public long max_coins(long[] nums) {
        long[] padded = new long[nums.Length + 2];
        padded[0] = 1;
        padded[nums.Length + 1] = 1;
        for (int i = 0; i < nums.Length; i++) padded[i + 1] = nums[i];
        int n = padded.Length;
        long[,] dp = new long[n, n];
        for (int length = 2; length < n; length++) {
            for (int left = 0; left < n - length; left++) {
                int right = left + length;
                for (int k = left + 1; k < right; k++) {
                    long coins = padded[left] * padded[k] * padded[right]
                                 + dp[left, k] + dp[k, right];
                    if (coins > dp[left, right]) dp[left, right] = coins;
                }
            }
        }
        return dp[0, n - 1];
    }
}
