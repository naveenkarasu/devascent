class Solution {
    public long collectible_count(long[] rain_times, long[] dry_durations) {
        long maxGap = 0;
        for (int i = 0; i < rain_times.length - 1; i++) {
            long gap = rain_times[i + 1] - rain_times[i];
            if (gap > maxGap) maxGap = gap;
        }
        long count = 0;
        for (long d : dry_durations) {
            if (d <= maxGap) count++;
        }
        return count;
    }
}
