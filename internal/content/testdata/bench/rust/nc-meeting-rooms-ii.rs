fn min_meeting_rooms(intervals: Vec<Vec<i64>>) -> i64 {
    if intervals.is_empty() {
        return 0;
    }
    let mut starts: Vec<i64> = intervals.iter().map(|i| i[0]).collect();
    let mut ends: Vec<i64> = intervals.iter().map(|i| i[1]).collect();
    starts.sort();
    ends.sort();
    let mut rooms: i64 = 0;
    let mut best: i64 = 0;
    let mut s = 0usize;
    let mut e = 0usize;
    while s < starts.len() {
        if starts[s] < ends[e] {
            rooms += 1;
            s += 1;
            if rooms > best {
                best = rooms;
            }
        } else {
            rooms -= 1;
            e += 1;
        }
    }
    best
}
