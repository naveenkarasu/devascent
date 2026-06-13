public class Solution {
    public string[][] solve(string[][] board) {
        if (board == null || board.Length == 0) return board;
        int m = board.Length, n = board[0].Length;
        bool[][] safe = new bool[m][];
        for (int i = 0; i < m; i++) safe[i] = new bool[n];

        for (int i = 0; i < m; i++) {
            Dfs(board, safe, i, 0, m, n);
            Dfs(board, safe, i, n-1, m, n);
        }
        for (int j = 0; j < n; j++) {
            Dfs(board, safe, 0, j, m, n);
            Dfs(board, safe, m-1, j, m, n);
        }

        string[][] result = new string[m][];
        for (int i = 0; i < m; i++) {
            result[i] = new string[n];
            for (int j = 0; j < n; j++) {
                if (board[i][j] == "O" && !safe[i][j]) result[i][j] = "X";
                else result[i][j] = board[i][j];
            }
        }
        return result;
    }

    private void Dfs(string[][] board, bool[][] safe, int i, int j, int m, int n) {
        if (i < 0 || j < 0 || i >= m || j >= n) return;
        if (board[i][j] != "O" || safe[i][j]) return;
        safe[i][j] = true;
        Dfs(board, safe, i+1, j, m, n);
        Dfs(board, safe, i-1, j, m, n);
        Dfs(board, safe, i, j+1, m, n);
        Dfs(board, safe, i, j-1, m, n);
    }
}
