fn wildcard_match(s: String, p: String) -> bool {
    let sc: Vec<char> = s.chars().collect();
    let pc: Vec<char> = p.chars().collect();
    let m = sc.len();
    let n = pc.len();
    let mut dp = vec![vec![false; n + 1]; m + 1];
    dp[0][0] = true;
    for j in 1..=n {
        if pc[j - 1] == '*' {
            dp[0][j] = dp[0][j - 1];
        }
    }
    for i in 1..=m {
        for j in 1..=n {
            if pc[j - 1] == '*' {
                dp[i][j] = dp[i - 1][j] || dp[i][j - 1];
            } else if pc[j - 1] == '?' || pc[j - 1] == sc[i - 1] {
                dp[i][j] = dp[i - 1][j - 1];
            }
        }
    }
    dp[m][n]
}
