fn max_square_divisor_root(n: i64) -> i64 {
    let mut answer: i64 = 1;
    let mut remaining = n;
    let mut fact: i64 = 2;
    while fact * fact <= n {
        let mut count = 0;
        while remaining % fact == 0 {
            count += 1;
            remaining /= fact;
            if count == 2 {
                answer *= fact;
                count = 0;
            }
        }
        fact += 1;
    }
    answer
}
