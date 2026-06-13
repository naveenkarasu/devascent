import java.util.*;

class Solution {
    public long majority_element(long[] nums) {
        long count = 0, candidate = 0;
        for (long n : nums) {
            if (count == 0) candidate = n;
            count += (n == candidate) ? 1 : -1;
        }
        long freq = 0;
        for (long n : nums) if (n == candidate) freq++;
        if (freq > nums.length / 2) return candidate;
        return -1;
    }
}
