fn counting_bits(n: i64) -> Vec<i64> {
    let n = n as usize;
    let mut dp = vec![0i64; n + 1];
    for i in 1..=n {
        dp[i] = dp[i >> 1] + (i & 1) as i64;
    }
    dp
}
