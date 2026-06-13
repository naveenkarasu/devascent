fn last_remaining(n: i64) -> i64 {
    if n == 1 {
        return 1;
    }
    2 * (1 + n / 2 - last_remaining(n / 2))
}
