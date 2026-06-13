public class Solution {
    public long max_profit(long[] prices) {
        long hold = long.MinValue / 2;
        long sold = 0;
        long rest = 0;
        foreach (var price in prices) {
            long prevSold = sold;
            sold = hold + price;
            hold = Math.Max(hold, rest - price);
            rest = Math.Max(rest, prevSold);
        }
        return Math.Max(sold, rest);
    }
}
