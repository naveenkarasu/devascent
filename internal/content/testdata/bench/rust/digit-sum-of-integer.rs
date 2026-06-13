fn digit_sum(n: i64) -> i64 {
    let mut n = n;
    let mut total: i64 = 0;
    while n > 0 {
        total += n % 10;
        n /= 10;
    }
    total
}
