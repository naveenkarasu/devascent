fn sum_squared_digits(n: i64) -> i64 {
    let mut n = n;
    let mut total: i64 = 0;
    while n > 0 {
        let d = n % 10;
        total += d * d;
        n /= 10;
    }
    total
}
