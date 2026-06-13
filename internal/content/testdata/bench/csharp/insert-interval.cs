public class Solution {
    public long[][] insert_interval(long[][] intervals, long[] newInterval) {
        var res = new List<long[]>();
        int i = 0, n = intervals.Length;
        long s = newInterval[0], e = newInterval[1];
        while (i < n && intervals[i][1] < s) { res.Add(intervals[i]); i++; }
        while (i < n && intervals[i][0] <= e) {
            s = Math.Min(s, intervals[i][0]);
            e = Math.Max(e, intervals[i][1]);
            i++;
        }
        res.Add(new long[] { s, e });
        while (i < n) { res.Add(intervals[i]); i++; }
        return res.ToArray();
    }
}
