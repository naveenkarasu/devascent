fn longest_ones_run(mut n: i64) -> i64 {
    let mut best = 0;
    let mut cur = 0;
    while n > 0 {
        if n & 1 == 1 {
            cur += 1;
            if cur > best {
                best = cur;
            }
        } else {
            cur = 0;
        }
        n >>= 1;
    }
    best
}
