import java.util.*;

class Solution {
    public long[][] rotate_matrix_cw(long[][] matrix) {
        int n = matrix.length;
        long[][] result = new long[n][n];
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                result[i][j] = matrix[n - 1 - j][i];
            }
        }
        return result;
    }
}
