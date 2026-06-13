function max_profit(prices: number[]): number {
    let hold = -Infinity;
    let sold = 0;
    let rest = 0;
    for (const price of prices) {
        const prev_sold = sold;
        sold = hold + price;
        hold = Math.max(hold, rest - price);
        rest = Math.max(rest, prev_sold);
    }
    return Math.max(sold, rest);
}
