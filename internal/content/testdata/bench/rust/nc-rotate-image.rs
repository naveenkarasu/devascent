fn rotate(matrix: Vec<Vec<i64>>) -> Vec<Vec<i64>> {
    let mut matrix = matrix;
    let n = matrix.len();
    for i in 0..n {
        for j in (i + 1)..n {
            let tmp = matrix[i][j];
            matrix[i][j] = matrix[j][i];
            matrix[j][i] = tmp;
        }
    }
    for row in matrix.iter_mut() {
        row.reverse();
    }
    matrix
}
