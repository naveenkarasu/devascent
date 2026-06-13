fn dfs(board: &Vec<Vec<String>>, safe: &mut Vec<Vec<bool>>, i: i64, j: i64, m: i64, n: i64) {
    if i < 0 || j < 0 || i >= m || j >= n {
        return;
    }
    let (ui, uj) = (i as usize, j as usize);
    if board[ui][uj] != "O" || safe[ui][uj] {
        return;
    }
    safe[ui][uj] = true;
    dfs(board, safe, i + 1, j, m, n);
    dfs(board, safe, i - 1, j, m, n);
    dfs(board, safe, i, j + 1, m, n);
    dfs(board, safe, i, j - 1, m, n);
}

fn solve(board: Vec<Vec<String>>) -> Vec<Vec<String>> {
    if board.is_empty() || board[0].is_empty() {
        return board;
    }
    let m = board.len() as i64;
    let n = board[0].len() as i64;
    let mut safe = vec![vec![false; board[0].len()]; board.len()];
    for i in 0..m {
        dfs(&board, &mut safe, i, 0, m, n);
        dfs(&board, &mut safe, i, n - 1, m, n);
    }
    for j in 0..n {
        dfs(&board, &mut safe, 0, j, m, n);
        dfs(&board, &mut safe, m - 1, j, m, n);
    }
    let mut result: Vec<Vec<String>> = Vec::new();
    for i in 0..(m as usize) {
        let mut row: Vec<String> = Vec::new();
        for j in 0..(n as usize) {
            if board[i][j] == "O" && !safe[i][j] {
                row.push("X".to_string());
            } else {
                row.push(board[i][j].clone());
            }
        }
        result.push(row);
    }
    result
}
