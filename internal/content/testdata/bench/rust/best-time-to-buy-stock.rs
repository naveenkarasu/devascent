fn max_profit(prices: Vec<i64>) -> i64 {
    let mut min_price = prices[0];
    let mut profit: i64 = 0;
    for &p in &prices[1..] {
        profit = profit.max(p - min_price);
        min_price = min_price.min(p);
    }
    profit
}
