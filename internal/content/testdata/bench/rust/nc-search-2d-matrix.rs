fn search_matrix(matrix: Vec<Vec<i64>>, target: i64) -> bool {
    let m = matrix.len() as i64;
    let n = matrix[0].len() as i64;
    let mut lo: i64 = 0;
    let mut hi: i64 = m * n - 1;
    while lo <= hi {
        let mid = (lo + hi) / 2;
        let val = matrix[(mid / n) as usize][(mid % n) as usize];
        if val == target {
            return true;
        } else if val < target {
            lo = mid + 1;
        } else {
            hi = mid - 1;
        }
    }
    false
}
