fn update_pair(a: i64, b: i64) -> Vec<i64> {
    vec![a + b, (a - b).abs()]
}
