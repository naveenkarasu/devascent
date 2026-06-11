fn palindrome_primes_in_range(a: i64, b: i64) -> Vec<i64> {
    fn isqrt(n: i64) -> i64 {
        if n < 0 {
            return 0;
        }
        let mut r = (n as f64).sqrt() as i64;
        while r > 0 && r * r > n {
            r -= 1;
        }
        while (r + 1) * (r + 1) <= n {
            r += 1;
        }
        r
    }
    fn is_prime(n: i64) -> bool {
        if n < 2 {
            return false;
        }
        let mut i = 2;
        let lim = isqrt(n);
        while i <= lim {
            if n % i == 0 {
                return false;
            }
            i += 1;
        }
        true
    }
    fn is_palindrome(n: i64) -> bool {
        let s: Vec<char> = n.to_string().chars().collect();
        let r: Vec<char> = s.iter().rev().cloned().collect();
        s == r
    }
    let mut result: Vec<i64> = Vec::new();
    for i in a..=b {
        if is_prime(i) && is_palindrome(i) {
            result.push(i);
        }
    }
    result
}
