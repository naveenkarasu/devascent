fn num_distinct(s: String, t: String) -> i64 {
    let sb = s.as_bytes();
    let tb = t.as_bytes();
    let m = sb.len();
    let n = tb.len();
    let mut dp = vec![0i64; n + 1];
    dp[0] = 1;
    for i in 1..=m {
        for j in (1..=n).rev() {
            if sb[i - 1] == tb[j - 1] {
                dp[j] += dp[j - 1];
            }
        }
    }
    dp[n]
}
