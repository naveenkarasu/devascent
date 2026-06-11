public class Solution {
    public long min_meeting_rooms(long[][] intervals) {
        if (intervals == null || intervals.Length == 0) return 0;
        var starts = intervals.Select(i => i[0]).OrderBy(x => x).ToArray();
        var ends = intervals.Select(i => i[1]).OrderBy(x => x).ToArray();
        long rooms = 0, best = 0;
        int s = 0, e = 0;
        while (s < starts.Length) {
            if (starts[s] < ends[e]) {
                rooms++;
                s++;
                best = Math.Max(best, rooms);
            } else {
                rooms--;
                e++;
            }
        }
        return best;
    }
}
