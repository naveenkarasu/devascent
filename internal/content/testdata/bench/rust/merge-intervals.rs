fn merge_intervals(intervals: Vec<Vec<i64>>) -> Vec<Vec<i64>> {
    let mut sorted = intervals.clone();
    sorted.sort();
    let mut res: Vec<Vec<i64>> = Vec::new();
    for iv in sorted {
        let s = iv[0];
        let e = iv[1];
        if let Some(last) = res.last_mut() {
            if s <= last[1] {
                if e > last[1] {
                    last[1] = e;
                }
                continue;
            }
        }
        res.push(vec![s, e]);
    }
    res
}
