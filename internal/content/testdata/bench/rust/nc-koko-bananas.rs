fn min_eating_speed(piles: Vec<i64>, h: i64) -> i64 {
    let mut lo: i64 = 1;
    let mut hi: i64 = *piles.iter().max().unwrap();
    let mut ans = hi;
    while lo <= hi {
        let mid = (lo + hi) / 2;
        let mut hours: i64 = 0;
        for &p in &piles {
            hours += (p + mid - 1) / mid; // ceil(p/mid) for positives
        }
        if hours <= h {
            ans = mid;
            hi = mid - 1;
        } else {
            lo = mid + 1;
        }
    }
    ans
}
