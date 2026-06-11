public class Solution {
    public long max_area(long[] heights) {
        int i = 0, j = heights.Length - 1;
        long best = 0;
        while (i < j) {
            long area = Math.Min(heights[i], heights[j]) * (j - i);
            best = Math.Max(best, area);
            if (heights[i] < heights[j]) i++;
            else j--;
        }
        return best;
    }
}
