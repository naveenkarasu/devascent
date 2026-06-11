fn coin_change(coins: Vec<i64>, amount: i64) -> i64 {
    let inf = amount + 1;
    let amt = amount as usize;
    let mut dp: Vec<i64> = vec![inf; amt + 1];
    dp[0] = 0;
    for a in 1..=amt {
        for &c in &coins {
            if c <= a as i64 {
                let prev = dp[a - c as usize] + 1;
                if prev < dp[a] {
                    dp[a] = prev;
                }
            }
        }
    }
    if dp[amt] != inf {
        dp[amt]
    } else {
        -1
    }
}
