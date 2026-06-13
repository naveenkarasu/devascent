public class Solution {
    public long find_celebrity(bool[][] knows_matrix) {
        int n = knows_matrix.Length;
        int candidate = 0;
        for (int i = 1; i < n; i++) {
            if (knows_matrix[candidate][i]) {
                candidate = i;
            }
        }
        for (int i = 0; i < n; i++) {
            if (i != candidate) {
                if (knows_matrix[candidate][i] || !knows_matrix[i][candidate]) {
                    return -1;
                }
            }
        }
        return candidate;
    }
}
