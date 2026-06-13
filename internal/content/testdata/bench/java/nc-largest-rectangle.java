import java.util.*;
class Solution {
    public long largest_rectangle_area(long[] heights) {
        long maxArea = 0;
        Deque<long[]> stack = new ArrayDeque<>(); // [start_index, height]
        int n = heights.length;
        for (int i = 0; i < n; i++) {
            long start = i;
            while (!stack.isEmpty() && stack.peek()[1] > heights[i]) {
                long[] top = stack.pop();
                long idx = top[0], height = top[1];
                maxArea = Math.max(maxArea, height * (i - idx));
                start = idx;
            }
            stack.push(new long[]{start, heights[i]});
        }
        for (long[] entry : stack) {
            maxArea = Math.max(maxArea, entry[1] * (n - entry[0]));
        }
        return maxArea;
    }
}
