fn clear_full_rows(board: Vec<Vec<String>>, empty: String) -> Vec<Vec<String>> {
    let cols = if !board.is_empty() { board[0].len() } else { 0 };
    let total = board.len();
    let surviving: Vec<Vec<String>> = board
        .into_iter()
        .filter(|row| row.iter().any(|c| *c == empty))
        .collect();
    let cleared_count = total - surviving.len();
    let mut result: Vec<Vec<String>> = Vec::new();
    for _ in 0..cleared_count {
        result.push(vec![empty.clone(); cols]);
    }
    result.extend(surviving);
    result
}
