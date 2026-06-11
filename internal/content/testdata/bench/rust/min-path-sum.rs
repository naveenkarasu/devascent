fn min_path_sum(grid: Vec<Vec<i64>>) -> i64 {
    let m = grid.len();
    let n = grid[0].len();
    let mut dp = grid.clone();
    for i in 0..m {
        for j in 0..n {
            if i == 0 && j == 0 {
                continue;
            } else if i == 0 {
                dp[i][j] += dp[i][j - 1];
            } else if j == 0 {
                dp[i][j] += dp[i - 1][j];
            } else {
                dp[i][j] += dp[i - 1][j].min(dp[i][j - 1]);
            }
        }
    }
    dp[m - 1][n - 1]
}
