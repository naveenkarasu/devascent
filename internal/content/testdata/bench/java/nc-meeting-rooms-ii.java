import java.util.Arrays;

class Solution {
    public long min_meeting_rooms(long[][] intervals) {
        if (intervals == null || intervals.length == 0) return 0;
        int n = intervals.length;
        long[] starts = new long[n];
        long[] ends = new long[n];
        for (int i = 0; i < n; i++) {
            starts[i] = intervals[i][0];
            ends[i] = intervals[i][1];
        }
        Arrays.sort(starts);
        Arrays.sort(ends);
        long rooms = 0, best = 0;
        int s = 0, e = 0;
        while (s < n) {
            if (starts[s] < ends[e]) {
                rooms++;
                s++;
                best = Math.max(best, rooms);
            } else {
                rooms--;
                e++;
            }
        }
        return best;
    }
}
