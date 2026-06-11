fn min_distance(word1: String, word2: String) -> i64 {
    let a = word1.as_bytes();
    let b = word2.as_bytes();
    let m = a.len();
    let n = b.len();
    let mut dp: Vec<i64> = (0..=n as i64).collect();
    for i in 1..=m {
        let mut prev = dp[0];
        dp[0] = i as i64;
        for j in 1..=n {
            let temp = dp[j];
            if a[i - 1] == b[j - 1] {
                dp[j] = prev;
            } else {
                let m1 = prev.min(dp[j]).min(dp[j - 1]);
                dp[j] = 1 + m1;
            }
            prev = temp;
        }
    }
    dp[n]
}
