fn is_match(s: String, p: String) -> bool {
    let sb = s.as_bytes();
    let pb = p.as_bytes();
    let m = sb.len();
    let n = pb.len();
    let mut dp = vec![vec![false; n + 1]; m + 1];
    dp[0][0] = true;
    for j in 2..=n {
        if pb[j - 1] == b'*' {
            dp[0][j] = dp[0][j - 2];
        }
    }
    for i in 1..=m {
        for j in 1..=n {
            if pb[j - 1] == b'*' {
                dp[i][j] = dp[i][j - 2];
                if pb[j - 2] == sb[i - 1] || pb[j - 2] == b'.' {
                    dp[i][j] = dp[i][j] || dp[i - 1][j];
                }
            } else if pb[j - 1] == sb[i - 1] || pb[j - 1] == b'.' {
                dp[i][j] = dp[i - 1][j - 1];
            }
        }
    }
    dp[m][n]
}
