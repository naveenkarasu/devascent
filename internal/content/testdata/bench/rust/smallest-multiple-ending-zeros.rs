fn gcd(a: i64, b: i64) -> i64 {
    let mut a = a;
    let mut b = b;
    while b != 0 {
        let t = b;
        b = a % b;
        a = t;
    }
    a
}

fn smallest_trailing_zero_multiple(a: i64, b: i64) -> i64 {
    let mut power: i64 = 1;
    for _ in 0..b {
        power *= 10;
    }
    (a / gcd(a, power)) * power
}
