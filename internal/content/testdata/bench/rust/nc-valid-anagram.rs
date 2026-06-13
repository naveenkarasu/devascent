use std::collections::HashMap;

fn is_anagram(s: String, t: String) -> bool {
    if s.len() != t.len() {
        return false;
    }
    let mut counts: HashMap<char, i64> = HashMap::new();
    for c in s.chars() {
        *counts.entry(c).or_insert(0) += 1;
    }
    for c in t.chars() {
        let entry = counts.entry(c).or_insert(0);
        if *entry == 0 {
            return false;
        }
        *entry -= 1;
    }
    true
}
