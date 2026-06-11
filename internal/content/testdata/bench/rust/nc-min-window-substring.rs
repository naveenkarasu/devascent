use std::collections::HashMap;

fn min_window(s: String, t: String) -> String {
    if t.is_empty() || s.is_empty() {
        return String::new();
    }
    let sc: Vec<char> = s.chars().collect();
    let tc: Vec<char> = t.chars().collect();
    let mut need: HashMap<char, i64> = HashMap::new();
    for &ch in tc.iter() {
        *need.entry(ch).or_insert(0) += 1;
    }
    let mut have: HashMap<char, i64> = HashMap::new();
    let mut formed: usize = 0;
    let required = need.len();
    let mut left: usize = 0;
    let mut best_len: i64 = i64::MAX;
    let mut best_start: usize = 0;
    for right in 0..sc.len() {
        let ch = sc[right];
        *have.entry(ch).or_insert(0) += 1;
        if need.contains_key(&ch) && have[&ch] == need[&ch] {
            formed += 1;
        }
        while formed == required {
            let cur_len = (right - left + 1) as i64;
            if cur_len < best_len {
                best_len = cur_len;
                best_start = left;
            }
            let left_ch = sc[left];
            if let Some(v) = have.get_mut(&left_ch) {
                *v -= 1;
            }
            if need.contains_key(&left_ch) && have[&left_ch] < need[&left_ch] {
                formed -= 1;
            }
            left += 1;
        }
    }
    if best_len == i64::MAX {
        return String::new();
    }
    sc[best_start..best_start + best_len as usize].iter().collect()
}
