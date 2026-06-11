public class Solution {
    public long max_subarray(long[] nums) {
        long current = nums[0];
        long best = nums[0];
        for (int i = 1; i < nums.Length; i++) {
            current = Math.Max(nums[i], current + nums[i]);
            best = Math.Max(best, current);
        }
        return best;
    }
}
