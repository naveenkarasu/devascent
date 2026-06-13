public class Solution {
    public long find_single(long[] nums) {
        long result = 0;
        foreach (long n in nums) {
            result ^= n;
        }
        return result;
    }
}
