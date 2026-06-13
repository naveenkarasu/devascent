public class Solution {
    public long maximize_score(long[] arr, long k) {
        long best = 0;
        foreach (var x in arr) {
            long diff = x - k;
            if (diff < 0) diff = -diff;
            long score = (2 * k) / (1 + 2 * diff);
            if (score > best) best = score;
        }
        return best;
    }
}
