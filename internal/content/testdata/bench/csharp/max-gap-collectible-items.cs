public class Solution {
    public long collectible_count(long[] rain_times, long[] dry_durations) {
        long maxGap = 0;
        for (int i = 0; i < rain_times.Length - 1; i++) {
            long gap = rain_times[i + 1] - rain_times[i];
            if (gap > maxGap) maxGap = gap;
        }
        long count = 0;
        foreach (long d in dry_durations) {
            if (d <= maxGap) count++;
        }
        return count;
    }
}
