using System.Linq;

public class Solution {
    public long kth_largest(long[] nums, long k) {
        long[] sorted = nums.OrderByDescending(x => x).ToArray();
        return sorted[(int)k - 1];
    }
}
