import java.util.*;

class Solution {
    public long max_profit(long[] prices) {
        long minPrice = prices[0];
        long profit = 0;
        for (int i = 1; i < prices.length; i++) {
            profit = Math.max(profit, prices[i] - minPrice);
            minPrice = Math.min(minPrice, prices[i]);
        }
        return profit;
    }
}
