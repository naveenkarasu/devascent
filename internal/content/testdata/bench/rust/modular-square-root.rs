fn modular_sqrt(n: i64, m: i64) -> i64 {
    let mut x = 0i64;
    while x < m {
        if (x * x) % m == n {
            return x;
        }
        x += 1;
    }
    -1
}
