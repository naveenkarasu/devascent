fn next_permutation(nums: Vec<i64>) -> Vec<i64> {
    let mut nums = nums;
    let n = nums.len() as i64;
    let mut i = n - 2;
    while i >= 0 && nums[i as usize] >= nums[(i + 1) as usize] {
        i -= 1;
    }
    if i != -1 {
        let mut j = n - 1;
        while nums[j as usize] <= nums[i as usize] {
            j -= 1;
        }
        nums.swap(i as usize, j as usize);
    }
    let mut left = (i + 1) as usize;
    let mut right = (n - 1) as usize;
    while left < right {
        nums.swap(left, right);
        left += 1;
        right -= 1;
    }
    nums
}
