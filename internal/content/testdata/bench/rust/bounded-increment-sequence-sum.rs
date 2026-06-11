fn bounded_sequence_sum(n: i64, b: i64, x: i64, y: i64) -> i64 {
    let n = n as usize;
    let mut a = vec![0i64; n + 1];
    for i in 1..=n {
        if a[i - 1] + x <= b {
            a[i] = a[i - 1] + x;
        } else {
            a[i] = a[i - 1] - y;
        }
    }
    a.iter().sum()
}
