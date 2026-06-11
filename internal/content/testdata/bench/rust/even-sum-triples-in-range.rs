fn count_even_sum_triples(arr: Vec<i64>, l: i64, r: i64) -> i64 {
    let lo = (l - 1) as usize;
    let hi = r as usize;
    let sub = &arr[lo..hi];
    let e = sub.iter().filter(|&&x| x % 2 == 0).count() as i64;
    let o = sub.len() as i64 - e;
    (e * (e - 1) * (e - 2)) / 6 + (o * (o - 1) / 2) * e
}
