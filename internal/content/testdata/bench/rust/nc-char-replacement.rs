use std::collections::HashMap;

fn character_replacement(s: String, k: i64) -> i64 {
    let chars: Vec<char> = s.chars().collect();
    let mut counts: HashMap<char, i64> = HashMap::new();
    let mut left: usize = 0;
    let mut max_count: i64 = 0;
    let mut best: i64 = 0;
    for right in 0..chars.len() {
        let ch = chars[right];
        let entry = counts.entry(ch).or_insert(0);
        *entry += 1;
        if *entry > max_count {
            max_count = *entry;
        }
        while (right as i64 - left as i64 + 1) - max_count > k {
            let lc = chars[left];
            if let Some(v) = counts.get_mut(&lc) {
                *v -= 1;
            }
            left += 1;
        }
        let window = right as i64 - left as i64 + 1;
        if window > best {
            best = window;
        }
    }
    best
}
