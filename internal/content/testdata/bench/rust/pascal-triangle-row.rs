fn pascal_row(n: i64) -> Vec<i64> {
    let mut result = Vec::new();
    let mut c: i64 = 1;
    for k in 0..n {
        result.push(c);
        // c = C(n-1, k+1) = C(n-1, k) * (n-1-k) / (k+1)
        c = c * (n - 1 - k) / (k + 1);
    }
    result
}
