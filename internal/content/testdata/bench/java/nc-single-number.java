class Solution {
    public long single_number(long[] nums) {
        long result = 0;
        for (long n : nums) result ^= n;
        return result;
    }
}
