fn reverse_integer(n: i64) -> i64 {
    let mut n = n;
    let mut reversed_n: i64 = 0;
    while n > 0 {
        reversed_n = reversed_n * 10 + n % 10;
        n /= 10;
    }
    reversed_n
}
