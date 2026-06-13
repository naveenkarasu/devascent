public class Solution {
    public long total_n_queens(long n) {
        int[] count = {0};
        var cols = new HashSet<int>();
        var diag1 = new HashSet<int>();
        var diag2 = new HashSet<int>();
        Backtrack(0, (int)n, cols, diag1, diag2, count);
        return count[0];
    }

    private void Backtrack(int row, int n, HashSet<int> cols, HashSet<int> diag1, HashSet<int> diag2, int[] count) {
        if (row == n) { count[0]++; return; }
        for (int col = 0; col < n; col++) {
            if (cols.Contains(col) || diag1.Contains(row - col) || diag2.Contains(row + col)) continue;
            cols.Add(col); diag1.Add(row - col); diag2.Add(row + col);
            Backtrack(row + 1, n, cols, diag1, diag2, count);
            cols.Remove(col); diag1.Remove(row - col); diag2.Remove(row + col);
        }
    }
}
