use std::collections::HashSet;

fn backtrack(
    row: i64,
    n: i64,
    cols: &mut HashSet<i64>,
    diag1: &mut HashSet<i64>,
    diag2: &mut HashSet<i64>,
    count: &mut i64,
) {
    if row == n {
        *count += 1;
        return;
    }
    for col in 0..n {
        if cols.contains(&col) || diag1.contains(&(row - col)) || diag2.contains(&(row + col)) {
            continue;
        }
        cols.insert(col);
        diag1.insert(row - col);
        diag2.insert(row + col);
        backtrack(row + 1, n, cols, diag1, diag2, count);
        cols.remove(&col);
        diag1.remove(&(row - col));
        diag2.remove(&(row + col));
    }
}

fn total_n_queens(n: i64) -> i64 {
    let mut count = 0i64;
    let mut cols: HashSet<i64> = HashSet::new();
    let mut diag1: HashSet<i64> = HashSet::new();
    let mut diag2: HashSet<i64> = HashSet::new();
    backtrack(0, n, &mut cols, &mut diag1, &mut diag2, &mut count);
    count
}
