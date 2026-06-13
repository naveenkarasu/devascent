fn my_pow(x: i64, n: i64) -> i64 {
    if n == 0 {
        return 1;
    }
    let half = my_pow(x, n / 2);
    if n % 2 == 0 {
        half * half
    } else {
        half * half * x
    }
}
