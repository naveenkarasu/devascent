import java.util.*;

class Solution {
    public long find_cheapest_price(long n, long[][] flights, long src, long dst, long k) {
        long INF = Long.MAX_VALUE / 2;
        long[] prices = new long[(int) n];
        Arrays.fill(prices, INF);
        prices[(int) src] = 0;
        for (int iter = 0; iter <= k; iter++) {
            long[] tmp = prices.clone();
            for (long[] flight : flights) {
                int u = (int) flight[0], v = (int) flight[1];
                long w = flight[2];
                if (prices[u] != INF && prices[u] + w < tmp[v]) {
                    tmp[v] = prices[u] + w;
                }
            }
            prices = tmp;
        }
        return prices[(int) dst] == INF ? -1 : prices[(int) dst];
    }
}
