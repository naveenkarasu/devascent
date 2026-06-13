fn is_perfect_square(n: i64) -> bool {
    if n < 0 {
        return false;
    }
    let mut r = (n as f64).sqrt() as i64;
    while r > 0 && r * r > n {
        r -= 1;
    }
    while (r + 1) * (r + 1) <= n {
        r += 1;
    }
    r * r == n
}
