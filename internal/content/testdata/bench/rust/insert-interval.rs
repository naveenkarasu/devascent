fn insert_interval(intervals: Vec<Vec<i64>>, new: Vec<i64>) -> Vec<Vec<i64>> {
    let mut res: Vec<Vec<i64>> = Vec::new();
    let n = intervals.len();
    let mut i = 0usize;
    let mut s = new[0];
    let mut e = new[1];
    while i < n && intervals[i][1] < s {
        res.push(intervals[i].clone());
        i += 1;
    }
    while i < n && intervals[i][0] <= e {
        s = s.min(intervals[i][0]);
        e = e.max(intervals[i][1]);
        i += 1;
    }
    res.push(vec![s, e]);
    while i < n {
        res.push(intervals[i].clone());
        i += 1;
    }
    res
}
