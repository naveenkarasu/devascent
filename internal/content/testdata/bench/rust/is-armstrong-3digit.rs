fn is_armstrong_3digit(n: i64) -> bool {
    if n < 100 || n > 999 {
        return false;
    }
    let mut s: i64 = 0;
    let mut temp = n;
    while temp > 0 {
        let r = temp % 10;
        s += r * r * r;
        temp /= 10;
    }
    s == n
}
