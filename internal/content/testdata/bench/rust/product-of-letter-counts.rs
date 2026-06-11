use std::collections::HashMap;

fn product_of_counts(s: String) -> i64 {
    const MOD: i64 = 1_000_000_007;
    let mut freq: HashMap<char, i64> = HashMap::new();
    for ch in s.chars() {
        *freq.entry(ch).or_insert(0) += 1;
    }
    let mut result: i64 = 1;
    for cnt in freq.values() {
        result = (result * cnt) % MOD;
    }
    result
}
