class Solution {
    public boolean subset_sum_equals_target(long[] nums, long target) {
        boolean[] dp = new boolean[(int) target + 1];
        dp[0] = true;
        for (long num : nums) {
            for (long j = target; j >= num; j--) {
                dp[(int) j] = dp[(int) j] || dp[(int) (j - num)];
            }
        }
        return dp[(int) target];
    }
}
