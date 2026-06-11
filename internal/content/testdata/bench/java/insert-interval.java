import java.util.*;

class Solution {
    public long[][] insert_interval(long[][] intervals, long[] newInterval) {
        List<long[]> res = new ArrayList<>();
        int i = 0, n = intervals.length;
        long s = newInterval[0], e = newInterval[1];
        while (i < n && intervals[i][1] < s) {
            res.add(intervals[i++]);
        }
        while (i < n && intervals[i][0] <= e) {
            s = Math.min(s, intervals[i][0]);
            e = Math.max(e, intervals[i][1]);
            i++;
        }
        res.add(new long[]{s, e});
        while (i < n) res.add(intervals[i++]);
        return res.toArray(new long[0][]);
    }
}
