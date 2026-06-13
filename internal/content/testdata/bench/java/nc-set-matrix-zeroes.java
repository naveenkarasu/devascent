import java.util.*;
class Solution {
    public long[][] set_zeroes(long[][] matrix) {
        Set<Integer> zeroRows = new HashSet<>();
        Set<Integer> zeroCols = new HashSet<>();
        int rows = matrix.length, cols = matrix[0].length;
        for (int r = 0; r < rows; r++)
            for (int c = 0; c < cols; c++)
                if (matrix[r][c] == 0) { zeroRows.add(r); zeroCols.add(c); }
        for (int r = 0; r < rows; r++)
            for (int c = 0; c < cols; c++)
                if (zeroRows.contains(r) || zeroCols.contains(c))
                    matrix[r][c] = 0;
        return matrix;
    }
}
