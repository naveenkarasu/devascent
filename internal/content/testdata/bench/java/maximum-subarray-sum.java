import java.util.*;

class Solution {
    public long max_subarray_sum(long[] nums) {
        long best = nums[0];
        long current = nums[0];
        for (int i = 1; i < nums.length; i++) {
            current = Math.max(nums[i], current + nums[i]);
            best = Math.max(best, current);
        }
        return best;
    }
}
