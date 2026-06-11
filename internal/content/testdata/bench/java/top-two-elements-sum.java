import java.util.*;
class Solution {
    public long top_two_sum(long[] nums) {
        long first = 0, second = 0;
        for (long x : nums) {
            if (x >= first) {
                second = first;
                first = x;
            } else if (x > second) {
                second = x;
            }
        }
        return first + second;
    }
}
