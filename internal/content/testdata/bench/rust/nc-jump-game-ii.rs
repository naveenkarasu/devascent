fn jump(nums: Vec<i64>) -> i64 {
    let mut jumps = 0i64;
    let mut cur_end = 0i64;
    let mut farthest = 0i64;
    let len = nums.len();
    if len == 0 {
        return 0;
    }
    for i in 0..(len - 1) {
        let reach = i as i64 + nums[i];
        if reach > farthest {
            farthest = reach;
        }
        if i as i64 == cur_end {
            jumps += 1;
            cur_end = farthest;
        }
    }
    jumps
}
