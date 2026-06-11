fn find_cheapest_price(n: i64, flights: Vec<Vec<i64>>, src: i64, dst: i64, k: i64) -> i64 {
    const INF: i64 = i64::MAX;
    let n = n as usize;
    let mut prices = vec![INF; n];
    prices[src as usize] = 0;
    for _ in 0..(k + 1) {
        let mut tmp = prices.clone();
        for f in &flights {
            let (u, v, w) = (f[0] as usize, f[1] as usize, f[2]);
            if prices[u] != INF && prices[u] + w < tmp[v] {
                tmp[v] = prices[u] + w;
            }
        }
        prices = tmp;
    }
    if prices[dst as usize] != INF {
        prices[dst as usize]
    } else {
        -1
    }
}
