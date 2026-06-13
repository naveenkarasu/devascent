public class Solution {
    public bool subset_sum_equals_target(long[] nums, long target) {
        bool[] dp = new bool[(int)target + 1];
        dp[0] = true;
        foreach (long num in nums) {
            for (long j = target; j >= num; j--) {
                dp[(int)j] = dp[(int)j] || dp[(int)(j - num)];
            }
        }
        return dp[(int)target];
    }
}
