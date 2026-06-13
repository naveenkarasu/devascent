fn set_union(a: Vec<i64>, b: Vec<i64>) -> Vec<i64> {
    let mut v: Vec<i64> = Vec::new();
    v.extend(a);
    v.extend(b);
    v.sort();
    v.dedup();
    v
}
