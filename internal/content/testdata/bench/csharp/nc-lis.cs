public class Solution {
    public long length_of_lis(long[] nums) {
        int n = nums.Length;
        long[] dp = new long[n];
        for (int i = 0; i < n; i++) dp[i] = 1;
        for (int i = 1; i < n; i++) {
            for (int j = 0; j < i; j++) {
                if (nums[j] < nums[i]) {
                    dp[i] = Math.Max(dp[i], dp[j] + 1);
                }
            }
        }
        long res = 0;
        foreach (var v in dp) res = Math.Max(res, v);
        return res;
    }
}
