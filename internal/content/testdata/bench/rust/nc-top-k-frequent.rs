fn top_k_frequent(nums: Vec<i64>, k: i64) -> Vec<i64> {
    use std::collections::HashMap;
    let mut counts: HashMap<i64, i64> = HashMap::new();
    for x in &nums {
        *counts.entry(*x).or_insert(0) += 1;
    }
    let mut keys: Vec<i64> = counts.keys().cloned().collect();
    // sort by (-count, value)
    keys.sort_by(|a, b| {
        let ca = counts[a];
        let cb = counts[b];
        cb.cmp(&ca).then(a.cmp(b))
    });
    let k = k as usize;
    let mut top: Vec<i64> = keys.into_iter().take(k).collect();
    top.sort();
    top
}
