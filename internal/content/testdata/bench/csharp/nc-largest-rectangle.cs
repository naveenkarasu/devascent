using System.Collections.Generic;

public class Solution {
    public long largest_rectangle_area(long[] heights) {
        long maxArea = 0;
        var stack = new Stack<long[]>(); // [start_index, height]
        int n = heights.Length;
        for (int i = 0; i < n; i++) {
            long start = i;
            while (stack.Count > 0 && stack.Peek()[1] > heights[i]) {
                long[] top = stack.Pop();
                long idx = top[0], height = top[1];
                maxArea = Math.Max(maxArea, height * (i - idx));
                start = idx;
            }
            stack.Push(new long[] { start, heights[i] });
        }
        foreach (long[] entry in stack) {
            maxArea = Math.Max(maxArea, entry[1] * (n - entry[0]));
        }
        return maxArea;
    }
}
