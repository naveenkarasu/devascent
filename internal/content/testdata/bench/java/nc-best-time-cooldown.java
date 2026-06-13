class Solution {
    public long max_profit(long[] prices) {
        long hold = Long.MIN_VALUE / 2;
        long sold = 0;
        long rest = 0;
        for (long price : prices) {
            long prevSold = sold;
            sold = hold + price;
            hold = Math.max(hold, rest - price);
            rest = Math.max(rest, prevSold);
        }
        return Math.max(sold, rest);
    }
}
