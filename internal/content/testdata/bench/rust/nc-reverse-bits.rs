fn reverse_bits(n: i64) -> i64 {
    let mut n = n;
    let mut result: i64 = 0;
    for _ in 0..32 {
        result = (result << 1) | (n & 1);
        n >>= 1;
    }
    result
}
