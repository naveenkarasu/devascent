func collectible_count(rain_times []int, dry_durations []int) int {
	maxGap := 0
	for i := 0; i < len(rain_times)-1; i++ {
		gap := rain_times[i+1] - rain_times[i]
		if gap > maxGap {
			maxGap = gap
		}
	}
	count := 0
	for _, d := range dry_durations {
		if d <= maxGap {
			count++
		}
	}
	return count
}
