function max_profit(prices: number[]): number {
    let buy = prices[0];
    let res = 0;
    for (let i = 1; i < prices.length; i++) {
        res = Math.max(res, prices[i] - buy);
        buy = Math.min(buy, prices[i]);
    }
    return res;
}
