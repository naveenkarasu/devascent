import java.util.*;

class Solution {
    public long max_area(long[] heights) {
        int i = 0, j = heights.length - 1;
        long best = 0;
        while (i < j) {
            long area = Math.min(heights[i], heights[j]) * (j - i);
            best = Math.max(best, area);
            if (heights[i] < heights[j]) i++;
            else j--;
        }
        return best;
    }
}
