public class Solution {
    public long majority_element(long[] nums) {
        long count = 0;
        long candidate = 0;
        foreach (long n in nums) {
            if (count == 0) candidate = n;
            count += (n == candidate) ? 1 : -1;
        }
        long cnt = 0;
        foreach (long n in nums) if (n == candidate) cnt++;
        if (cnt > nums.Length / 2) return candidate;
        return -1;
    }
}
