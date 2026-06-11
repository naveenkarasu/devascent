fn edit_distance(word1: String, word2: String) -> i64 {
    let w1: Vec<char> = word1.chars().collect();
    let w2: Vec<char> = word2.chars().collect();
    let m = w1.len();
    let n = w2.len();
    let mut dp: Vec<i64> = (0..=(n as i64)).collect();
    for i in 1..=m {
        let mut prev = dp[0];
        dp[0] = i as i64;
        for j in 1..=n {
            let temp = dp[j];
            if w1[i - 1] == w2[j - 1] {
                dp[j] = prev;
            } else {
                dp[j] = 1 + prev.min(dp[j - 1]).min(dp[j]);
            }
            prev = temp;
        }
    }
    dp[n]
}
