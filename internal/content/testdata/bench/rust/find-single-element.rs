fn find_single(nums: Vec<i64>) -> i64 {
    let mut result: i64 = 0;
    for n in nums {
        result ^= n;
    }
    result
}
