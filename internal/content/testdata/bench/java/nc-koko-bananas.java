class Solution {
    public long min_eating_speed(long[] piles, long h) {
        long lo = 1, hi = 0;
        for (long p : piles) hi = Math.max(hi, p);
        long ans = hi;
        while (lo <= hi) {
            long mid = (lo + hi) / 2;
            long hours = 0;
            for (long p : piles) hours += (p + mid - 1) / mid;
            if (hours <= h) {
                ans = mid;
                hi = mid - 1;
            } else {
                lo = mid + 1;
            }
        }
        return ans;
    }
}
