fn subset_sum_to_k(nums: Vec<i64>, k: i64) -> bool {
    let k = k as usize;
    let n = nums.len();
    let mut prev = vec![false; k + 1];
    prev[0] = true;
    if (nums[0] as usize) <= k {
        prev[nums[0] as usize] = true;
    }
    for i in 1..n {
        let mut curr = vec![false; k + 1];
        curr[0] = true;
        for j in 1..=k {
            let not_take = prev[j];
            let take = if (nums[i] as usize) <= j {
                prev[j - nums[i] as usize]
            } else {
                false
            };
            curr[j] = take || not_take;
        }
        prev = curr;
    }
    prev[k]
}
