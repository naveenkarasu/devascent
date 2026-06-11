fn max_coins(nums: Vec<i64>) -> i64 {
    let mut padded = Vec::with_capacity(nums.len() + 2);
    padded.push(1i64);
    padded.extend(nums.iter().cloned());
    padded.push(1i64);
    let n = padded.len();
    let mut dp = vec![vec![0i64; n]; n];
    for length in 2..n {
        for left in 0..(n - length) {
            let right = left + length;
            for k in (left + 1)..right {
                let coins = padded[left] * padded[k] * padded[right]
                    + dp[left][k]
                    + dp[k][right];
                if coins > dp[left][right] {
                    dp[left][right] = coins;
                }
            }
        }
    }
    dp[0][n - 1]
}
