fn modular_exponentiation(x: i64, n: i64, m: i64) -> i64 {
    let mut res: i64 = 1;
    let mut x = x % m;
    let mut n = n;
    while n > 0 {
        if n % 2 != 0 {
            res = (res * x) % m;
            n -= 1;
        } else {
            x = (x * x) % m;
            n /= 2;
        }
    }
    res
}
