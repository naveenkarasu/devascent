fn rotate_matrix_cw(matrix: Vec<Vec<i64>>) -> Vec<Vec<i64>> {
    let n = matrix.len();
    let mut result: Vec<Vec<i64>> = Vec::with_capacity(n);
    for i in 0..n {
        let mut row: Vec<i64> = Vec::with_capacity(n);
        for j in 0..n {
            row.push(matrix[n - 1 - j][i]);
        }
        result.push(row);
    }
    result
}
