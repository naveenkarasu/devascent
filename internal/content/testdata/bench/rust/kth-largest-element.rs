fn kth_largest(nums: Vec<i64>, k: i64) -> i64 {
    let mut nums = nums;
    nums.sort_by(|a, b| b.cmp(a));
    nums[(k - 1) as usize]
}
