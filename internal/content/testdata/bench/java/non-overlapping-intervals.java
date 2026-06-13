import java.util.*;

class Solution {
    public long erase_overlap_intervals(long[][] intervals) {
        if (intervals == null || intervals.length == 0) return 0;
        Arrays.sort(intervals, (a, b) -> Long.compare(a[1], b[1]));
        long end = intervals[0][1];
        long count = 0;
        for (int i = 1; i < intervals.length; i++) {
            if (intervals[i][0] < end) {
                count++;
            } else {
                end = intervals[i][1];
            }
        }
        return count;
    }
}
