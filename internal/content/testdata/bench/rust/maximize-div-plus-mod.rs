fn max_div_plus_mod(l: i64, r: i64, a: i64) -> i64 {
    let mut best = r / a + r % a;
    let multiple = (r / a) * a;
    let candidate = multiple - 1;
    if l <= candidate && candidate <= r {
        best = best.max(candidate / a + candidate % a);
    }
    best
}
