class Solution {
    public boolean exist(String[][] board, String word) {
        int rows = board.length, cols = board[0].length;
        boolean[][] visited = new boolean[rows][cols];
        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                if (dfs(board, word, r, c, 0, visited, rows, cols)) return true;
            }
        }
        return false;
    }

    private boolean dfs(String[][] board, String word, int r, int c, int index, boolean[][] visited, int rows, int cols) {
        if (index == word.length()) return true;
        if (r < 0 || r >= rows || c < 0 || c >= cols) return false;
        if (visited[r][c] || !board[r][c].equals(String.valueOf(word.charAt(index)))) return false;
        visited[r][c] = true;
        boolean found = dfs(board, word, r+1, c, index+1, visited, rows, cols)
                     || dfs(board, word, r-1, c, index+1, visited, rows, cols)
                     || dfs(board, word, r, c+1, index+1, visited, rows, cols)
                     || dfs(board, word, r, c-1, index+1, visited, rows, cols);
        visited[r][c] = false;
        return found;
    }
}
