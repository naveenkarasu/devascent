fn sort_012(nums: Vec<i64>) -> Vec<i64> {
    let mut nums = nums;
    let mut low = 0i64;
    let mut mid = 0i64;
    let mut high = nums.len() as i64 - 1;
    while mid <= high {
        if nums[mid as usize] == 0 {
            nums.swap(low as usize, mid as usize);
            low += 1;
            mid += 1;
        } else if nums[mid as usize] == 1 {
            mid += 1;
        } else {
            nums.swap(mid as usize, high as usize);
            high -= 1;
        }
    }
    nums
}
