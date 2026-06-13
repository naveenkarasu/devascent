fn unique_paths(m: i64, n: i64) -> i64 {
    let n = n as usize;
    let mut dp = vec![1i64; n];
    for _ in 1..m {
        for j in 1..n {
            dp[j] += dp[j - 1];
        }
    }
    dp[n - 1]
}
