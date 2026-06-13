function find_cheapest_price(n, flights, src, dst, k) {
  const INF = Infinity;
  let prices = new Array(n).fill(INF);
  prices[src] = 0;

  for (let i = 0; i <= k; i++) {
    const tmp = prices.slice();
    for (const [u, v, w] of flights) {
      if (prices[u] !== INF && prices[u] + w < tmp[v]) {
        tmp[v] = prices[u] + w;
      }
    }
    prices = tmp;
  }

  return prices[dst] === INF ? -1 : prices[dst];
}
