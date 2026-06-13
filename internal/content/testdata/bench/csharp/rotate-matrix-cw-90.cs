public class Solution {
    public long[][] rotate_matrix_cw(long[][] matrix) {
        int n = matrix.Length;
        long[][] result = new long[n][];
        for (int i = 0; i < n; i++) {
            result[i] = new long[n];
            for (int j = 0; j < n; j++) {
                result[i][j] = matrix[n - 1 - j][i];
            }
        }
        return result;
    }
}
