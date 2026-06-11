fn length_of_lis(nums: Vec<i64>) -> i64 {
    let n = nums.len();
    if n == 0 {
        return 0;
    }
    let mut dp = vec![1i64; n];
    for i in 1..n {
        for j in 0..i {
            if nums[j] < nums[i] && dp[j] + 1 > dp[i] {
                dp[i] = dp[j] + 1;
            }
        }
    }
    *dp.iter().max().unwrap()
}
