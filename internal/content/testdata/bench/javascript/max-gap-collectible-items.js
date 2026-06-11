function collectible_count(rain_times, dry_durations) {
    let maxGap = 0;
    for (let i = 0; i < rain_times.length - 1; i++) {
        let gap = rain_times[i + 1] - rain_times[i];
        if (gap > maxGap) maxGap = gap;
    }
    let count = 0;
    for (let d of dry_durations) {
        if (d <= maxGap) count++;
    }
    return count;
}
