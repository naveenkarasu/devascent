fn max_subarray_sum(nums: Vec<i64>) -> i64 {
    let mut best = nums[0];
    let mut current = nums[0];
    for &x in &nums[1..] {
        current = x.max(current + x);
        best = best.max(current);
    }
    best
}
