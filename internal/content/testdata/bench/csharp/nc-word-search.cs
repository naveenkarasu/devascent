public class Solution {
    public bool exist(string[][] board, string word) {
        int rows = board.Length, cols = board[0].Length;
        bool[,] visited = new bool[rows, cols];
        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                if (Dfs(board, word, r, c, 0, visited, rows, cols)) return true;
            }
        }
        return false;
    }

    private bool Dfs(string[][] board, string word, int r, int c, int index, bool[,] visited, int rows, int cols) {
        if (index == word.Length) return true;
        if (r < 0 || r >= rows || c < 0 || c >= cols) return false;
        if (visited[r, c] || board[r][c] != word[index].ToString()) return false;
        visited[r, c] = true;
        bool found = Dfs(board, word, r+1, c, index+1, visited, rows, cols)
                  || Dfs(board, word, r-1, c, index+1, visited, rows, cols)
                  || Dfs(board, word, r, c+1, index+1, visited, rows, cols)
                  || Dfs(board, word, r, c-1, index+1, visited, rows, cols);
        visited[r, c] = false;
        return found;
    }
}
