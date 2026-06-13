function max_profit(prices: number[]): number {
    let minPrice = prices[0];
    let profit = 0;
    for (let i = 1; i < prices.length; i++) {
        profit = Math.max(profit, prices[i] - minPrice);
        minPrice = Math.min(minPrice, prices[i]);
    }
    return profit;
}
