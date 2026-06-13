fn is_fibonacci(n: i64) -> bool {
    if n <= 0 {
        return false;
    }
    let mut a: i64 = 1;
    let mut b: i64 = 1;
    while a < n {
        let next = a + b;
        a = b;
        b = next;
    }
    a == n
}
