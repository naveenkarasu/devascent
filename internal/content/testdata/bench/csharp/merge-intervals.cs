public class Solution {
    public long[][] merge_intervals(long[][] intervals) {
        var sorted = intervals.OrderBy(x => x[0]).ToArray();
        var res = new List<long[]>();
        foreach (var interval in sorted) {
            long s = interval[0], e = interval[1];
            if (res.Count > 0 && s <= res[res.Count - 1][1]) {
                res[res.Count - 1][1] = Math.Max(res[res.Count - 1][1], e);
            } else {
                res.Add(new long[] { s, e });
            }
        }
        return res.ToArray();
    }
}
