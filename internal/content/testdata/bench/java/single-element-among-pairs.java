class Solution {
    public long find_single(long[] arr) {
        long result = 0;
        for (long x : arr) result ^= x;
        return result;
    }
}
