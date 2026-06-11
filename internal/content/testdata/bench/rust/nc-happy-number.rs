use std::collections::HashSet;

fn is_happy(n: i64) -> bool {
    let mut n = n;
    let mut seen: HashSet<i64> = HashSet::new();
    while n != 1 {
        if seen.contains(&n) {
            return false;
        }
        seen.insert(n);
        let mut total = 0i64;
        let mut m = n;
        while m != 0 {
            let digit = m % 10;
            m = m / 10;
            total += digit * digit;
        }
        n = total;
    }
    true
}
