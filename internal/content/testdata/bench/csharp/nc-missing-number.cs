public class Solution {
    public long missing_number(long[] nums) {
        long n = nums.Length;
        long expected = n * (n + 1) / 2;
        long sum = 0;
        foreach (long x in nums) sum += x;
        return expected - sum;
    }
}
