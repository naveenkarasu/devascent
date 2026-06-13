fn reverse_integer(x: i64) -> i64 {
    let int_min: i64 = -(1i64 << 31);
    let int_max: i64 = (1i64 << 31) - 1;
    let sign: i64 = if x < 0 { -1 } else { 1 };
    let mut n = x.abs();
    let mut rev: i64 = 0;
    while n != 0 {
        rev = rev * 10 + n % 10;
        n /= 10;
    }
    let result = sign * rev;
    if result < int_min || result > int_max {
        return 0;
    }
    result
}
