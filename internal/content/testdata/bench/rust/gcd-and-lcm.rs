fn gcd_lcm(a: i64, b: i64) -> Vec<i64> {
    fn gcd(mut x: i64, mut y: i64) -> i64 {
        while y != 0 {
            let t = x % y;
            x = y;
            y = t;
        }
        x
    }
    let g = gcd(a, b);
    let l = a * b / g;
    vec![g, l]
}
