fn reverse_integer(x: i64) -> i64 {
    let sign: i64 = if x < 0 { -1 } else { 1 };
    let mut x = x.abs();
    let mut rev: i64 = 0;
    while x > 0 {
        rev = rev * 10 + x % 10;
        x /= 10;
    }
    rev *= sign;
    if rev > 2147483647 || rev < -2147483648 {
        return 0;
    }
    rev
}
