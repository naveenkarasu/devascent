use std::collections::HashMap;

fn length_of_longest_substring(s: String) -> i64 {
    let chars: Vec<char> = s.chars().collect();
    let mut last_seen: HashMap<char, i64> = HashMap::new();
    let mut left: i64 = 0;
    let mut best: i64 = 0;
    for (right, &ch) in chars.iter().enumerate() {
        let right = right as i64;
        if let Some(&prev) = last_seen.get(&ch) {
            if prev >= left {
                left = prev + 1;
            }
        }
        last_seen.insert(ch, right);
        if right - left + 1 > best {
            best = right - left + 1;
        }
    }
    best
}
