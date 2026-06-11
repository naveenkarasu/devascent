import java.util.*;
class Solution {
    public long count_pairs_same_mod(long[] nums, long divisor) {
        Map<Long, Long> freq = new HashMap<>();
        for (long x : nums) {
            long rem = x % divisor;
            freq.put(rem, freq.getOrDefault(rem, 0L) + 1);
        }
        long pairs = 0;
        for (long cnt : freq.values()) {
            pairs += cnt * (cnt - 1) / 2;
        }
        return pairs;
    }
}
