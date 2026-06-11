fn change(amount: i64, coins: Vec<i64>) -> i64 {
    let amt = amount as usize;
    let mut dp = vec![0i64; amt + 1];
    dp[0] = 1;
    for &coin in coins.iter() {
        let c = coin as usize;
        for x in c..=amt {
            dp[x] += dp[x - c];
        }
    }
    dp[amt]
}
