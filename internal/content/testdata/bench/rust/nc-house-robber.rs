fn rob(nums: Vec<i64>) -> i64 {
    if nums.is_empty() {
        return 0;
    }
    if nums.len() == 1 {
        return nums[0];
    }
    let mut prev2 = nums[0];
    let mut prev1 = nums[0].max(nums[1]);
    for i in 2..nums.len() {
        let curr = prev1.max(prev2 + nums[i]);
        prev2 = prev1;
        prev1 = curr;
    }
    prev1
}
