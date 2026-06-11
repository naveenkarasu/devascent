fn lexical_order(n: i64) -> Vec<i64> {
    let mut result: Vec<i64> = Vec::new();
    let mut cur: i64 = 1;
    while (result.len() as i64) < n {
        result.push(cur);
        if cur * 10 <= n {
            cur *= 10;
        } else {
            while cur % 10 == 9 || cur >= n {
                cur /= 10;
            }
            cur += 1;
        }
    }
    result
}
