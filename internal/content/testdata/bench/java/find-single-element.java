import java.util.*;

class Solution {
    public long find_single(long[] nums) {
        long result = 0;
        for (long n : nums) {
            result ^= n;
        }
        return result;
    }
}
