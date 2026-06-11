public class Solution {
    public long single_number(long[] nums) {
        long result = 0;
        foreach (long n in nums) result ^= n;
        return result;
    }
}
