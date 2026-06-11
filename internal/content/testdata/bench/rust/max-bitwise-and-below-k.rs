fn max_and_below_k(n: i64, k: i64) -> i64 {
    if ((k - 1) | k) <= n {
        k - 1
    } else {
        k - 2
    }
}
