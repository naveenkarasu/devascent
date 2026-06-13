fn dfs(
    board: &Vec<Vec<String>>,
    word: &Vec<char>,
    visited: &mut Vec<Vec<bool>>,
    r: i64,
    c: i64,
    index: usize,
    rows: i64,
    cols: i64,
) -> bool {
    if index == word.len() {
        return true;
    }
    if r < 0 || r >= rows || c < 0 || c >= cols {
        return false;
    }
    let (ur, uc) = (r as usize, c as usize);
    if visited[ur][uc] || board[ur][uc] != word[index].to_string() {
        return false;
    }
    visited[ur][uc] = true;
    let found = dfs(board, word, visited, r + 1, c, index + 1, rows, cols)
        || dfs(board, word, visited, r - 1, c, index + 1, rows, cols)
        || dfs(board, word, visited, r, c + 1, index + 1, rows, cols)
        || dfs(board, word, visited, r, c - 1, index + 1, rows, cols);
    visited[ur][uc] = false;
    found
}

fn exist(board: Vec<Vec<String>>, word: String) -> bool {
    let rows = board.len();
    let cols = board[0].len();
    let word_chars: Vec<char> = word.chars().collect();
    let mut visited = vec![vec![false; cols]; rows];
    for r in 0..rows {
        for c in 0..cols {
            if dfs(
                &board,
                &word_chars,
                &mut visited,
                r as i64,
                c as i64,
                0,
                rows as i64,
                cols as i64,
            ) {
                return true;
            }
        }
    }
    false
}
