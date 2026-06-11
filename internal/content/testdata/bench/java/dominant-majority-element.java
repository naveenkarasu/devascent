import java.util.*;

class Solution {
    public long find_majority_element(long[] nums) {
        long candidate = 0;
        int count = 0;
        for (long x : nums) {
            if (count == 0) candidate = x;
            count += (x == candidate) ? 1 : -1;
        }
        int freq = 0;
        for (long x : nums) if (x == candidate) freq++;
        if (freq > nums.length / 2) return candidate;
        return -1;
    }
}
