fn max_hourglass_sum(grid: Vec<Vec<i64>>) -> i64 {
    let mut best = i64::MIN;
    for i in 0..4 {
        for j in 0..4 {
            let s = grid[i][j] + grid[i][j + 1] + grid[i][j + 2]
                + grid[i + 1][j + 1]
                + grid[i + 2][j] + grid[i + 2][j + 1] + grid[i + 2][j + 2];
            if s > best {
                best = s;
            }
        }
    }
    best
}
