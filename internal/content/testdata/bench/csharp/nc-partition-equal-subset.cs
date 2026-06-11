using System.Collections.Generic;

public class Solution {
    public bool can_partition(long[] nums) {
        long total = 0;
        foreach (var n in nums) total += n;
        if (total % 2 != 0) return false;
        long target = total / 2;
        var dp = new HashSet<long> { 0 };
        foreach (var n in nums) {
            var next = new HashSet<long>(dp);
            foreach (var s in dp) next.Add(s + n);
            dp = next;
            if (dp.Contains(target)) return true;
        }
        return false;
    }
}
