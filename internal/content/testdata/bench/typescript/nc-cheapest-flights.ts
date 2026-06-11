function find_cheapest_price(n: number, flights: number[][], src: number, dst: number, k: number): number {
    const INF = Number.MAX_SAFE_INTEGER;
    let prices: number[] = new Array(n).fill(INF);
    prices[src] = 0;
    for (let step = 0; step <= k; step++) {
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
