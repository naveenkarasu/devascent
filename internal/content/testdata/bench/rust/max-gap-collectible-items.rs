fn collectible_count(rain_times: Vec<i64>, dry_durations: Vec<i64>) -> i64 {
    let mut max_gap: i64 = 0;
    for i in 0..rain_times.len().saturating_sub(1) {
        let gap = rain_times[i + 1] - rain_times[i];
        if gap > max_gap {
            max_gap = gap;
        }
    }
    dry_durations.into_iter().filter(|&d| d <= max_gap).count() as i64
}
