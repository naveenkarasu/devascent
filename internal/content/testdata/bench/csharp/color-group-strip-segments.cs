using System.Collections.Generic;

public class Solution {
    public long count_segments(long[][] groups, long[] strip) {
        var colorToGroup = new Dictionary<long, int>();
        for (int idx = 0; idx < groups.Length; idx++) {
            foreach (long color in groups[idx]) {
                colorToGroup[color] = idx;
            }
        }
        if (strip.Length == 0) return 0;
        long segments = 1;
        for (int i = 1; i < strip.Length; i++) {
            int g1 = colorToGroup.GetValueOrDefault(strip[i - 1], -1);
            int g2 = colorToGroup.GetValueOrDefault(strip[i], -1);
            if (g1 != g2) segments++;
        }
        return segments;
    }
}
