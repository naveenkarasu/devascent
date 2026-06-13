import java.util.*;

class Solution {
    public long peak_occupancy(long initial, long[] deltas) {
        long current = initial;
        long peak = initial;
        for (long d : deltas) {
            current += d;
            if (current > peak) peak = current;
        }
        return peak;
    }
}
