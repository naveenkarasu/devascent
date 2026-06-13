public class Solution {
    public long find_single(long[] arr) {
        long result = 0;
        foreach (long x in arr) result ^= x;
        return result;
    }
}
