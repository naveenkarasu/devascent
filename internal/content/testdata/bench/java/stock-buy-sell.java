import java.util.*;

class Solution {
    public long max_profit(long[] prices) {
        long buy = prices[0], res = 0;
        for (int i = 1; i < prices.length; i++) {
            res = Math.max(res, prices[i] - buy);
            buy = Math.min(buy, prices[i]);
        }
        return res;
    }
}
