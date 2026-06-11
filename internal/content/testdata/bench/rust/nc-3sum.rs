fn three_sum(nums: Vec<i64>) -> Vec<Vec<i64>> {
    let mut nums = nums;
    nums.sort();
    let mut res: Vec<Vec<i64>> = Vec::new();
    let n = nums.len();
    for i in 0..n {
        if i > 0 && nums[i] == nums[i - 1] {
            continue;
        }
        let mut lo = i + 1;
        let mut hi = n.saturating_sub(1);
        while lo < hi {
            let total = nums[i] + nums[lo] + nums[hi];
            if total < 0 {
                lo += 1;
            } else if total > 0 {
                hi -= 1;
            } else {
                res.push(vec![nums[i], nums[lo], nums[hi]]);
                lo += 1;
                hi -= 1;
                while lo < hi && nums[lo] == nums[lo - 1] {
                    lo += 1;
                }
                while lo < hi && nums[hi] == nums[hi + 1] {
                    hi -= 1;
                }
            }
        }
    }
    res
}
