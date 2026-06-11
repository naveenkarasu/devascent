import java.util.*;
class Solution {
    public long[] match_colored_houses(long[] left_colors, long[] right_colors) {
        Map<Long, Long> colorToPos = new HashMap<>();
        for (int i = 0; i < right_colors.length; i++) {
            colorToPos.put(right_colors[i], (long)(i + 1));
        }
        long[] result = new long[left_colors.length];
        for (int i = 0; i < left_colors.length; i++) {
            result[i] = colorToPos.get(left_colors[i]);
        }
        return result;
    }
}
