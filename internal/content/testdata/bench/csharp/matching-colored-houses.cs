using System.Collections.Generic;

public class Solution {
    public long[] match_colored_houses(long[] left_colors, long[] right_colors) {
        var colorToPos = new Dictionary<long, long>();
        for (int i = 0; i < right_colors.Length; i++) {
            colorToPos[right_colors[i]] = i + 1;
        }
        var result = new long[left_colors.Length];
        for (int i = 0; i < left_colors.Length; i++) {
            result[i] = colorToPos[left_colors[i]];
        }
        return result;
    }
}
