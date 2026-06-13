fn is_factorion(n: i64) -> bool {
    fn factorial(d: i64) -> i64 {
        let mut f: i64 = 1;
        let mut i = 2;
        while i <= d {
            f *= i;
            i += 1;
        }
        f
    }
    let original = n;
    let mut n = n;
    let mut s: i64 = 0;
    while n > 0 {
        let d = n % 10;
        s += factorial(d);
        n /= 10;
    }
    s == original
}
