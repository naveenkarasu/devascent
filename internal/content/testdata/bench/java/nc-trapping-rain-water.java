import java.util.*;

class Solution {
    public long trap(long[] heights) {
        if (heights.length == 0) return 0;
        int i = 0, j = heights.length - 1;
        long leftMax = heights[i], rightMax = heights[j];
        long total = 0;
        while (i < j) {
            if (leftMax < rightMax) {
                i++;
                leftMax = Math.max(leftMax, heights[i]);
                total += leftMax - heights[i];
            } else {
                j--;
                rightMax = Math.max(rightMax, heights[j]);
                total += rightMax - heights[j];
            }
        }
        return total;
    }
}
