fn steps_to_overtake(a: i64, b: i64) -> i64 {
    let mut a = a;
    let mut b = b;
    let mut steps = 0i64;
    while a <= b {
        a *= 3;
        b *= 2;
        steps += 1;
    }
    steps
}
