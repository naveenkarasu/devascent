import java.util.*;

class Solution {
    public boolean can_attend_meetings(long[][] intervals) {
        Arrays.sort(intervals, (a, b) -> Long.compare(a[0], b[0]));
        for (int i = 1; i < intervals.length; i++) {
            if (intervals[i][0] < intervals[i-1][1]) return false;
        }
        return true;
    }
}
