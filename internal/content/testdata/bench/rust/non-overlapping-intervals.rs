fn erase_overlap_intervals(intervals: Vec<Vec<i64>>) -> i64 {
    if intervals.is_empty() {
        return 0;
    }
    let mut sorted = intervals.clone();
    sorted.sort_by(|a, b| a[1].cmp(&b[1]));
    let mut end = sorted[0][1];
    let mut count: i64 = 0;
    for i in 1..sorted.len() {
        if sorted[i][0] < end {
            count += 1;
        } else {
            end = sorted[i][1];
        }
    }
    count
}
