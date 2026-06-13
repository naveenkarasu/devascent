fn free_vertices(n: i64, m: i64) -> i64 {
    // Find smallest k such that k*(k-1)/2 >= m.
    let mut k: i64 = 0;
    while k * (k - 1) / 2 < m {
        k += 1;
    }
    n - k
}
