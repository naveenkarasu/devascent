import java.util.*;

class Solution {
    public long find_target_sum_ways(long[] nums, long target) {
        Map<Long, Long> dp = new HashMap<>();
        dp.put(0L, 1L);
        for (long num : nums) {
            Map<Long, Long> nextDp = new HashMap<>();
            for (Map.Entry<Long, Long> entry : dp.entrySet()) {
                long s = entry.getKey(), cnt = entry.getValue();
                nextDp.merge(s + num, cnt, Long::sum);
                nextDp.merge(s - num, cnt, Long::sum);
            }
            dp = nextDp;
        }
        return dp.getOrDefault(target, 0L);
    }
}
