import java.util.*;

class Solution {
    public long count_segments(long[][] groups, long[] strip) {
        Map<Long, Integer> colorToGroup = new HashMap<>();
        for (int idx = 0; idx < groups.length; idx++) {
            for (long color : groups[idx]) {
                colorToGroup.put(color, idx);
            }
        }
        if (strip.length == 0) return 0;
        long segments = 1;
        for (int i = 1; i < strip.length; i++) {
            Integer g1 = colorToGroup.get(strip[i - 1]);
            Integer g2 = colorToGroup.get(strip[i]);
            if (!Objects.equals(g1, g2)) {
                segments++;
            }
        }
        return segments;
    }
}
