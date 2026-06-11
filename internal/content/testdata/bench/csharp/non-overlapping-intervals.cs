public class Solution {
    public long erase_overlap_intervals(long[][] intervals) {
        if (intervals == null || intervals.Length == 0) return 0;
        var sorted = intervals.OrderBy(x => x[1]).ToArray();
        long end = sorted[0][1];
        long count = 0;
        for (int i = 1; i < sorted.Length; i++) {
            if (sorted[i][0] < end) {
                count++;
            } else {
                end = sorted[i][1];
            }
        }
        return count;
    }
}
