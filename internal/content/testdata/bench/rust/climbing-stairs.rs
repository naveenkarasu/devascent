fn climb_stairs(n: i64) -> i64 {
    let mut a: i64 = 1;
    let mut b: i64 = 1;
    for _ in 0..n {
        let next = a + b;
        a = b;
        b = next;
    }
    a
}
