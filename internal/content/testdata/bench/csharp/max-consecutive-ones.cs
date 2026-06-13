public class Solution {
    public long max_consecutive_ones(long[] nums) {
        long maxRun = 0, current = 0;
        foreach (long x in nums) {
            if (x == 1) {
                current++;
                if (current > maxRun) maxRun = current;
            } else {
                current = 0;
            }
        }
        return maxRun;
    }
}
