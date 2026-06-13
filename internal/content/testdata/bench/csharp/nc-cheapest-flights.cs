public class Solution {
    public long find_cheapest_price(long n, long[][] flights, long src, long dst, long k) {
        long INF = long.MaxValue / 2;
        long[] prices = new long[(int)n];
        for (int i = 0; i < n; i++) prices[i] = INF;
        prices[(int)src] = 0;
        for (int iter = 0; iter <= k; iter++) {
            long[] tmp = (long[])prices.Clone();
            foreach (var flight in flights) {
                int u = (int)flight[0], v = (int)flight[1];
                long w = flight[2];
                if (prices[u] != INF && prices[u] + w < tmp[v]) {
                    tmp[v] = prices[u] + w;
                }
            }
            prices = tmp;
        }
        return prices[(int)dst] == INF ? -1 : prices[(int)dst];
    }
}
