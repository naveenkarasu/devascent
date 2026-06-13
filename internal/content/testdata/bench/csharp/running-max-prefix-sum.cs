public class Solution {
    public long peak_occupancy(long initial, long[] deltas) {
        long current = initial;
        long peak = initial;
        foreach (long d in deltas) {
            current += d;
            if (current > peak) peak = current;
        }
        return peak;
    }
}
