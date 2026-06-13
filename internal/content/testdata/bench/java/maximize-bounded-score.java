class Solution {
    public long maximize_score(long[] arr, long k) {
        long best = 0;
        for (long x : arr) {
            long score = (2 * k) / (1 + 2 * Math.abs(x - k));
            if (score > best) best = score;
        }
        return best;
    }
}
