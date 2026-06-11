fn maximize_score(arr: Vec<i64>, k: i64) -> i64 {
    let mut best: i64 = i64::MIN;
    for x in arr {
        let score = (2 * k).div_euclid(1 + 2 * (x - k).abs());
        if score > best {
            best = score;
        }
    }
    best
}
