public class Solution {
    public long max_profit(long[] prices) {
        long minPrice = prices[0];
        long profit = 0;
        for (int i = 1; i < prices.Length; i++) {
            profit = System.Math.Max(profit, prices[i] - minPrice);
            minPrice = System.Math.Min(minPrice, prices[i]);
        }
        return profit;
    }
}
