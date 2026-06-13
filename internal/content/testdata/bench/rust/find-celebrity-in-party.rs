fn find_celebrity(knows_matrix: Vec<Vec<bool>>) -> i64 {
    let n = knows_matrix.len();
    let mut candidate = 0usize;
    for i in 1..n {
        if knows_matrix[candidate][i] {
            candidate = i;
        }
    }
    for i in 0..n {
        if i != candidate {
            if knows_matrix[candidate][i] || !knows_matrix[i][candidate] {
                return -1;
            }
        }
    }
    candidate as i64
}
