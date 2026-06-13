import java.util.*;
class Solution {
    public long max_profit(long[] prices) {
        long minPrice = Long.MAX_VALUE, best = 0;
        for (long p : prices) {
            if (p < minPrice) minPrice = p;
            else if (p - minPrice > best) best = p - minPrice;
        }
        return best;
    }
}
