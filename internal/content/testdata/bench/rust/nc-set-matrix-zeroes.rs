use std::collections::HashSet;

fn set_zeroes(matrix: Vec<Vec<i64>>) -> Vec<Vec<i64>> {
    let mut matrix = matrix;
    let mut zero_rows: HashSet<usize> = HashSet::new();
    let mut zero_cols: HashSet<usize> = HashSet::new();
    let rows = matrix.len();
    let cols = if rows > 0 { matrix[0].len() } else { 0 };
    for r in 0..rows {
        for c in 0..cols {
            if matrix[r][c] == 0 {
                zero_rows.insert(r);
                zero_cols.insert(c);
            }
        }
    }
    for r in 0..rows {
        for c in 0..cols {
            if zero_rows.contains(&r) || zero_cols.contains(&c) {
                matrix[r][c] = 0;
            }
        }
    }
    matrix
}
