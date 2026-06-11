fn sum_of_factorials(n: i64) -> i64 {
    let mut total: i64 = 0;
    let mut fact: i64 = 1;
    for i in 1..=n {
        fact *= i;
        total += fact;
    }
    total
}
