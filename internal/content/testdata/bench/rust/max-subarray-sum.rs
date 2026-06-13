fn max_subarray(nums: Vec<i64>) -> i64 {
    let mut current = nums[0];
    let mut best = nums[0];
    for i in 1..nums.len() {
        current = nums[i].max(current + nums[i]);
        best = best.max(current);
    }
    best
}
