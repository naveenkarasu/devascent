function collectible_count(rain_times: number[], dry_durations: number[]): number {
    let max_gap = 0;
    for (let i = 0; i < rain_times.length - 1; i++) {
        const gap = rain_times[i + 1] - rain_times[i];
        if (gap > max_gap) max_gap = gap;
    }
    return dry_durations.filter(d => d <= max_gap).length;
}
