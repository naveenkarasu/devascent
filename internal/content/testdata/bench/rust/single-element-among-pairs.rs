fn find_single(arr: Vec<i64>) -> i64 {
    let mut result = 0;
    for x in arr {
        result ^= x;
    }
    result
}
