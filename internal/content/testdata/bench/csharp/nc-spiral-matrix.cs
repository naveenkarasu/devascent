using System.Collections.Generic;

public class Solution {
    public long[] spiral_order(long[][] matrix) {
        var result = new List<long>();
        int top = 0, bottom = matrix.Length - 1;
        int left = 0, right = matrix[0].Length - 1;
        while (top <= bottom && left <= right) {
            for (int col = left; col <= right; col++) result.Add(matrix[top][col]);
            top++;
            for (int row = top; row <= bottom; row++) result.Add(matrix[row][right]);
            right--;
            if (top <= bottom) {
                for (int col = right; col >= left; col--) result.Add(matrix[bottom][col]);
                bottom--;
            }
            if (left <= right) {
                for (int row = bottom; row >= top; row--) result.Add(matrix[row][left]);
                left++;
            }
        }
        return result.ToArray();
    }
}
