import java.util.*;

class Solution {
    public boolean can_partition(long[] nums) {
        long total = 0;
        for (long n : nums) total += n;
        if (total % 2 != 0) return false;
        long target = total / 2;
        Set<Long> dp = new HashSet<>();
        dp.add(0L);
        for (long n : nums) {
            Set<Long> next = new HashSet<>(dp);
            for (long s : dp) {
                next.add(s + n);
            }
            dp = next;
            if (dp.contains(target)) return true;
        }
        return false;
    }
}
