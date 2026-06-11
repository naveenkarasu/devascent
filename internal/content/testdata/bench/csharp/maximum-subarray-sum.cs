public class Solution {
    public long max_subarray_sum(long[] nums) {
        long best = nums[0];
        long current = nums[0];
        for (int i = 1; i < nums.Length; i++) {
            current = System.Math.Max(nums[i], current + nums[i]);
            best = System.Math.Max(best, current);
        }
        return best;
    }
}
