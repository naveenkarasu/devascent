fn max_profit(prices: Vec<i64>) -> i64 {
    let mut buy = prices[0];
    let mut res = 0i64;
    for i in 1..prices.len() {
        res = res.max(prices[i] - buy);
        buy = buy.min(prices[i]);
    }
    res
}
