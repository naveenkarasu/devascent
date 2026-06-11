import java.util.*;

class Solution {
    public long[] two_sum(long[] nums, long target) {
        Map<Long, Integer> seen = new HashMap<>();
        for (int i = 0; i < nums.length; i++) {
            long need = target - nums[i];
            if (seen.containsKey(need)) return new long[]{seen.get(need), i};
            seen.put(nums[i], i);
        }
        return new long[]{};
    }
}
