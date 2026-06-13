import java.util.*;
class Solution {
    public long[][] rotate(long[][] matrix) {
        int n = matrix.length;
        for (int i = 0; i < n; i++) {
            for (int j = i + 1; j < n; j++) {
                long tmp = matrix[i][j];
                matrix[i][j] = matrix[j][i];
                matrix[j][i] = tmp;
            }
        }
        for (int i = 0; i < n; i++) {
            for (int l = 0, r = n - 1; l < r; l++, r--) {
                long tmp = matrix[i][l];
                matrix[i][l] = matrix[i][r];
                matrix[i][r] = tmp;
            }
        }
        return matrix;
    }
}
