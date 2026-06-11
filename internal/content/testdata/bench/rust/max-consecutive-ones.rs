fn max_consecutive_ones(nums: Vec<i64>) -> i64 {
    let mut max_run = 0;
    let mut current = 0;
    for x in nums {
        if x == 1 {
            current += 1;
            if current > max_run {
                max_run = current;
            }
        } else {
            current = 0;
        }
    }
    max_run
}
