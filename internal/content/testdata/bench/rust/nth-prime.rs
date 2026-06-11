fn nth_prime(n: i64) -> i64 {
    let mut primes: Vec<i64> = Vec::new();
    let mut i: i64 = 2;
    while (primes.len() as i64) < n {
        let mut is_p = true;
        for &p in &primes {
            if i % p == 0 {
                is_p = false;
                break;
            }
        }
        if is_p {
            primes.push(i);
        }
        i += 1;
    }
    *primes.last().unwrap()
}
