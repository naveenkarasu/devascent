fn is_interleave(s1: String, s2: String, s3: String) -> bool {
    let a = s1.as_bytes();
    let b = s2.as_bytes();
    let c = s3.as_bytes();
    let m = a.len();
    let n = b.len();
    if m + n != c.len() {
        return false;
    }
    let mut dp = vec![false; n + 1];
    dp[0] = true;
    for j in 1..=n {
        dp[j] = dp[j - 1] && b[j - 1] == c[j - 1];
    }
    for i in 1..=m {
        dp[0] = dp[0] && a[i - 1] == c[i - 1];
        for j in 1..=n {
            dp[j] = (dp[j] && a[i - 1] == c[i + j - 1])
                || (dp[j - 1] && b[j - 1] == c[i + j - 1]);
        }
    }
    dp[n]
}
