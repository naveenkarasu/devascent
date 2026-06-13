class Solution {
    public long max_coins(long[] nums) {
        int origLen = nums.length;
        long[] arr = new long[origLen + 2];
        arr[0] = 1;
        arr[origLen + 1] = 1;
        for (int i = 0; i < origLen; i++) arr[i + 1] = nums[i];
        int n = arr.length;
        long[][] dp = new long[n][n];
        for (int length = 2; length < n; length++) {
            for (int left = 0; left < n - length; left++) {
                int right = left + length;
                for (int k = left + 1; k < right; k++) {
                    long coins = arr[left] * arr[k] * arr[right] + dp[left][k] + dp[k][right];
                    if (coins > dp[left][right]) dp[left][right] = coins;
                }
            }
        }
        return dp[0][n - 1];
    }
}
