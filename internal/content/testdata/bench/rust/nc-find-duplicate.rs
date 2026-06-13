fn find_duplicate(nums: Vec<i64>) -> i64 {
    let mut slow = nums[0] as usize;
    let mut fast = nums[nums[0] as usize] as usize;
    while slow != fast {
        slow = nums[slow] as usize;
        fast = nums[nums[fast] as usize] as usize;
    }
    let mut slow2 = 0usize;
    while slow != slow2 {
        slow = nums[slow] as usize;
        slow2 = nums[slow2] as usize;
    }
    slow as i64
}
