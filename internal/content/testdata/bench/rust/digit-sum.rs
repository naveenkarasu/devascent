fn digit_sum(n: i64) -> i64 {
    let mut n = n;
    let mut s: i64 = 0;
    while n > 0 {
        s += n % 10;
        n /= 10;
    }
    s
}
