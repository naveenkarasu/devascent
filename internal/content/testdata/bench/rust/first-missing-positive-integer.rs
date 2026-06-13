fn first_missing_positive(mut nums: Vec<i64>) -> i64 {
    let n = nums.len() as i64;
    let mut i = 0usize;
    while i < nums.len() {
        while nums[i] >= 1 && nums[i] <= n {
            let correct = (nums[i] - 1) as usize;
            if nums[correct] == nums[i] {
                break;
            }
            nums.swap(i, correct);
        }
        i += 1;
    }
    for i in 0..nums.len() {
        if nums[i] != (i as i64) + 1 {
            return (i as i64) + 1;
        }
    }
    n + 1
}
