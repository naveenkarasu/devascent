using System.Collections.Generic;

public class Solution {
    public long find_target_sum_ways(long[] nums, long target) {
        var dp = new Dictionary<long, long> { { 0, 1 } };
        foreach (var num in nums) {
            var nextDp = new Dictionary<long, long>();
            foreach (var kv in dp) {
                long s = kv.Key;
                long cnt = kv.Value;
                long addKey = s + num;
                long subKey = s - num;
                nextDp.TryGetValue(addKey, out long v1);
                nextDp[addKey] = v1 + cnt;
                nextDp.TryGetValue(subKey, out long v2);
                nextDp[subKey] = v2 + cnt;
            }
            dp = nextDp;
        }
        dp.TryGetValue(target, out long result);
        return result;
    }
}
