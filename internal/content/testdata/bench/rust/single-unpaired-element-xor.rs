fn find_unpaired(arr: Vec<i64>) -> i64 {
    let mut result: i64 = 0;
    for x in arr {
        result ^= x;
    }
    result
}
