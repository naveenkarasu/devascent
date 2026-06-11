function exist(board, word) {
  const rows = board.length, cols = board[0].length;
  const visited = new Set();

  function dfs(r, c, index) {
    if (index === word.length) return true;
    if (r < 0 || r >= rows || c < 0 || c >= cols) return false;
    if (visited.has(r * cols + c)) return false;
    if (board[r][c] !== word[index]) return false;
    visited.add(r * cols + c);
    const found = dfs(r+1, c, index+1) || dfs(r-1, c, index+1) ||
                  dfs(r, c+1, index+1) || dfs(r, c-1, index+1);
    visited.delete(r * cols + c);
    return found;
  }

  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      if (dfs(r, c, 0)) return true;
    }
  }
  return false;
}
