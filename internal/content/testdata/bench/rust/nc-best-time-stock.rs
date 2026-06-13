fn max_profit(prices: Vec<i64>) -> i64 {
    let mut min_price = i64::MAX;
    let mut best = 0i64;
    for &p in prices.iter() {
        if p < min_price {
            min_price = p;
        } else if p - min_price > best {
            best = p - min_price;
        }
    }
    best
}
