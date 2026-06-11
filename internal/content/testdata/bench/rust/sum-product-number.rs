fn is_sum_product_number(n: i64) -> bool {
    let mut digit_sum: i64 = 0;
    let mut digit_prod: i64 = 1;
    let mut temp = n;
    while temp > 0 {
        let d = temp % 10;
        digit_sum += d;
        digit_prod *= d;
        temp /= 10;
    }
    digit_sum * digit_prod == n
}
