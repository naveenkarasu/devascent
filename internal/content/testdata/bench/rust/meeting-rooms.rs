fn can_attend_meetings(intervals: Vec<Vec<i64>>) -> bool {
    let mut sorted = intervals.clone();
    sorted.sort();
    for i in 1..sorted.len() {
        if sorted[i][0] < sorted[i - 1][1] {
            return false;
        }
    }
    true
}
