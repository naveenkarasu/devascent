public class Solution {
    public long find_unpaired(long[] arr) {
        long result = 0;
        foreach (var x in arr) {
            result ^= x;
        }
        return result;
    }
}
