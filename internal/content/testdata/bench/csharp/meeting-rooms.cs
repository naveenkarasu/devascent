public class Solution {
    public bool can_attend_meetings(long[][] intervals) {
        var sorted = intervals.OrderBy(x => x[0]).ToArray();
        for (int i = 1; i < sorted.Length; i++) {
            if (sorted[i][0] < sorted[i - 1][1]) return false;
        }
        return true;
    }
}
