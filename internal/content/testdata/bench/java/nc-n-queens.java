import java.util.*;

class Solution {
    public long total_n_queens(long n) {
        int[] count = {0};
        Set<Integer> cols = new HashSet<>(), diag1 = new HashSet<>(), diag2 = new HashSet<>();
        backtrack(0, (int) n, cols, diag1, diag2, count);
        return count[0];
    }

    private void backtrack(int row, int n, Set<Integer> cols, Set<Integer> diag1, Set<Integer> diag2, int[] count) {
        if (row == n) { count[0]++; return; }
        for (int col = 0; col < n; col++) {
            if (cols.contains(col) || diag1.contains(row - col) || diag2.contains(row + col)) continue;
            cols.add(col); diag1.add(row - col); diag2.add(row + col);
            backtrack(row + 1, n, cols, diag1, diag2, count);
            cols.remove(col); diag1.remove(row - col); diag2.remove(row + col);
        }
    }
}
