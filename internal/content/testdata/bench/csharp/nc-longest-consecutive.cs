public class Solution {
    public long longest_consecutive(long[] nums) {
        var set = new HashSet<long>(nums);
        long best = 0;
        foreach (long n in set) {
            if (!set.Contains(n - 1)) {
                long length = 1;
                while (set.Contains(n + length)) length++;
                best = Math.Max(best, length);
            }
        }
        return best;
    }
}
