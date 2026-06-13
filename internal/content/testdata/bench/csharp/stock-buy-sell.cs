public class Solution {
    public long max_profit(long[] prices) {
        long buy = prices[0];
        long res = 0;
        for (int i = 1; i < prices.Length; i++) {
            res = Math.Max(res, prices[i] - buy);
            buy = Math.Min(buy, prices[i]);
        }
        return res;
    }
}
