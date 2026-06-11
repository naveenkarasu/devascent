public class Solution {
    public long find_majority_element(long[] nums) {
        long candidate = 0;
        long count = 0;
        foreach (long x in nums) {
            if (count == 0) candidate = x;
            count += (x == candidate) ? 1 : -1;
        }
        long freq = 0;
        foreach (long x in nums) if (x == candidate) freq++;
        if (freq > nums.Length / 2) return candidate;
        return -1;
    }
}
