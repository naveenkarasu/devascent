fn num_decodings(s: String) -> i64 {
    let b = s.as_bytes();
    let n = b.len();
    if n == 0 {
        return 0;
    }
    let mut dp = vec![0i64; n + 1];
    dp[0] = 1;
    dp[1] = if b[0] == b'0' { 0 } else { 1 };
    for i in 2..=n {
        let one = (b[i - 1] - b'0') as i64;
        let two = (b[i - 2] - b'0') as i64 * 10 + (b[i - 1] - b'0') as i64;
        if one != 0 {
            dp[i] += dp[i - 1];
        }
        if two >= 10 && two <= 26 {
            dp[i] += dp[i - 2];
        }
    }
    dp[n]
}
