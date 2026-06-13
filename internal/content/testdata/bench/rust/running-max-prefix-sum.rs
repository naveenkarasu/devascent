fn peak_occupancy(initial: i64, deltas: Vec<i64>) -> i64 {
    let mut current = initial;
    let mut peak = initial;
    for d in deltas {
        current += d;
        if current > peak {
            peak = current;
        }
    }
    peak
}
