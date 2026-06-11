class Solution {
    public long missing_number(long[] nums) {
        long n = nums.length;
        long expected = n * (n + 1) / 2;
        long sum = 0;
        for (long x : nums) sum += x;
        return expected - sum;
    }
}
