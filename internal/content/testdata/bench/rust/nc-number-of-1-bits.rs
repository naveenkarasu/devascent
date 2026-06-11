fn hamming_weight(n: i64) -> i64 {
    let mut n = n;
    let mut count: i64 = 0;
    while n != 0 {
        count += n & 1;
        n >>= 1;
    }
    count
}
