fn subset_sum_equals_target(nums: Vec<i64>, target: i64) -> bool {
    if target < 0 {
        return false;
    }
    let t = target as usize;
    let mut dp = vec![false; t + 1];
    dp[0] = true;
    for &num in &nums {
        if num < 0 {
            continue;
        }
        let num_u = num as usize;
        if num_u > t {
            continue;
        }
        for j in (num_u..=t).rev() {
            dp[j] = dp[j] || dp[j - num_u];
        }
    }
    dp[t]
}
