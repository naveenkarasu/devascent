fn spiral_order(matrix: Vec<Vec<i64>>) -> Vec<i64> {
    let mut result: Vec<i64> = Vec::new();
    if matrix.is_empty() || matrix[0].is_empty() {
        return result;
    }
    let mut top = 0i64;
    let mut bottom = matrix.len() as i64 - 1;
    let mut left = 0i64;
    let mut right = matrix[0].len() as i64 - 1;
    while top <= bottom && left <= right {
        let mut col = left;
        while col <= right {
            result.push(matrix[top as usize][col as usize]);
            col += 1;
        }
        top += 1;
        let mut row = top;
        while row <= bottom {
            result.push(matrix[row as usize][right as usize]);
            row += 1;
        }
        right -= 1;
        if top <= bottom {
            let mut col = right;
            while col >= left {
                result.push(matrix[bottom as usize][col as usize]);
                col -= 1;
            }
            bottom -= 1;
        }
        if left <= right {
            let mut row = bottom;
            while row >= top {
                result.push(matrix[row as usize][left as usize]);
                row -= 1;
            }
            left += 1;
        }
    }
    result
}
