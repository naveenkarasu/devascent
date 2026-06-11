fn days_breakdown(d: i64) -> Vec<i64> {
    let y = d / 365;
    let mut d = d % 365;
    let w = d / 7;
    d = d % 7;
    vec![y, w, d]
}
