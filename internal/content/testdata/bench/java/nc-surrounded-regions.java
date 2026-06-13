import java.util.*;

class Solution {
    public String[][] solve(String[][] board) {
        if (board == null || board.length == 0) return board;
        int m = board.length, n = board[0].length;
        boolean[][] safe = new boolean[m][n];

        for (int i = 0; i < m; i++) {
            dfs(board, safe, i, 0, m, n);
            dfs(board, safe, i, n-1, m, n);
        }
        for (int j = 0; j < n; j++) {
            dfs(board, safe, 0, j, m, n);
            dfs(board, safe, m-1, j, m, n);
        }

        String[][] result = new String[m][n];
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if ("O".equals(board[i][j]) && !safe[i][j]) {
                    result[i][j] = "X";
                } else {
                    result[i][j] = board[i][j];
                }
            }
        }
        return result;
    }

    private void dfs(String[][] board, boolean[][] safe, int i, int j, int m, int n) {
        if (i < 0 || j < 0 || i >= m || j >= n) return;
        if (!"O".equals(board[i][j]) || safe[i][j]) return;
        safe[i][j] = true;
        dfs(board, safe, i+1, j, m, n);
        dfs(board, safe, i-1, j, m, n);
        dfs(board, safe, i, j+1, m, n);
        dfs(board, safe, i, j-1, m, n);
    }
}
