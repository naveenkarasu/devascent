class Solution {
    public long length_of_lis(long[] nums) {
        long[] dp = new long[nums.length];
        java.util.Arrays.fill(dp, 1);
        for (int i = 1; i < nums.length; i++) {
            for (int j = 0; j < i; j++) {
                if (nums[j] < nums[i]) {
                    dp[i] = Math.max(dp[i], dp[j] + 1);
                }
            }
        }
        long max = 0;
        for (long v : dp) max = Math.max(max, v);
        return max;
    }
}
