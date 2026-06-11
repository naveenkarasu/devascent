fn count_modeq_pairs(n: i64, m: i64) -> i64 {
    let mut count = 0i64;
    let mut a = 1i64;
    while a <= n {
        let mut b = a + 1;
        while b <= n {
            if (m % a) % b == (m % b) % a {
                count += 1;
            }
            b += 1;
        }
        a += 1;
    }
    count
}
