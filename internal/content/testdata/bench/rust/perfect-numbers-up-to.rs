fn perfect_numbers_up_to(n: i64) -> Vec<i64> {
    fn isqrt(n: i64) -> i64 {
        if n < 0 {
            return 0;
        }
        let mut r = (n as f64).sqrt() as i64;
        while r > 0 && r * r > n {
            r -= 1;
        }
        while (r + 1) * (r + 1) <= n {
            r += 1;
        }
        r
    }
    let mut result: Vec<i64> = Vec::new();
    for i in 2..=n {
        let mut s: i64 = 1;
        let r = isqrt(i);
        let mut d = 2;
        while d <= r {
            if i % d == 0 {
                s += d + i / d;
                if d * d == i {
                    s -= d;
                }
            }
            d += 1;
        }
        if s == i {
            result.push(i);
        }
    }
    result
}
