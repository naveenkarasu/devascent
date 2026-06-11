import java.util.*;

class Solution {
    public long longest_consecutive(long[] nums) {
        Set<Long> set = new HashSet<>();
        for (long n : nums) set.add(n);
        long best = 0;
        for (long n : set) {
            if (!set.contains(n - 1)) {
                long length = 1;
                while (set.contains(n + length)) length++;
                best = Math.max(best, length);
            }
        }
        return best;
    }
}
