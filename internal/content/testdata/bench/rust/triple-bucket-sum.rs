fn triple_bucket_sum(counts: Vec<i64>) -> i64 {
    const MOD: i64 = 1_000_000_007;
    let mut ans: i64 = 0;
    let n = counts.len();
    if n < 3 {
        return 0;
    }
    for i in 0..n - 2 {
        for j in i + 1..n - 1 {
            for k in j + 1..n {
                ans = (ans + counts[i] * counts[j] % MOD * counts[k]) % MOD;
            }
        }
    }
    ans
}
