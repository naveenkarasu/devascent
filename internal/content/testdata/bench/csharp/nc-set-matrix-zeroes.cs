using System.Collections.Generic;

public class Solution {
    public long[][] set_zeroes(long[][] matrix) {
        var zeroRows = new HashSet<int>();
        var zeroCols = new HashSet<int>();
        int rows = matrix.Length, cols = matrix[0].Length;
        for (int r = 0; r < rows; r++)
            for (int c = 0; c < cols; c++)
                if (matrix[r][c] == 0) { zeroRows.Add(r); zeroCols.Add(c); }
        for (int r = 0; r < rows; r++)
            for (int c = 0; c < cols; c++)
                if (zeroRows.Contains(r) || zeroCols.Contains(c))
                    matrix[r][c] = 0;
        return matrix;
    }
}
