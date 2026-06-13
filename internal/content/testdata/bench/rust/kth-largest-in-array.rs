fn find_kth_largest(nums: Vec<i64>, k: i64) -> i64 {
    let mut sorted = nums.clone();
    sorted.sort_unstable();
    let idx = sorted.len() - k as usize;
    sorted[idx]
}
