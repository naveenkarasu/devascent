use std::collections::HashMap;

fn count_pairs_same_mod(nums: Vec<i64>, divisor: i64) -> i64 {
    let mut freq: HashMap<i64, i64> = HashMap::new();
    for x in &nums {
        let key = x.rem_euclid(divisor);
        *freq.entry(key).or_insert(0) += 1;
    }
    let mut pairs: i64 = 0;
    for cnt in freq.values() {
        pairs += cnt * (cnt - 1) / 2;
    }
    pairs
}
