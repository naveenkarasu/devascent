public class Solution {
    public long max_profit(long[] prices) {
        long minPrice = long.MaxValue, best = 0;
        foreach (long p in prices) {
            if (p < minPrice) minPrice = p;
            else if (p - minPrice > best) best = p - minPrice;
        }
        return best;
    }
}
