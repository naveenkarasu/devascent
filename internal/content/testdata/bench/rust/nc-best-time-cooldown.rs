fn max_profit(prices: Vec<i64>) -> i64 {
    let neg_inf = i64::MIN / 2;
    let mut hold = neg_inf;
    let mut sold = 0i64;
    let mut rest = 0i64;
    for &price in prices.iter() {
        let prev_sold = sold;
        sold = hold + price;
        hold = hold.max(rest - price);
        rest = rest.max(prev_sold);
    }
    sold.max(rest)
}
