use std::collections::HashMap;

fn check_inclusion(s1: String, s2: String) -> bool {
    let a: Vec<char> = s1.chars().collect();
    let b: Vec<char> = s2.chars().collect();
    if a.len() > b.len() {
        return false;
    }
    let mut need: HashMap<char, i64> = HashMap::new();
    for &ch in a.iter() {
        *need.entry(ch).or_insert(0) += 1;
    }
    let mut window: HashMap<char, i64> = HashMap::new();
    let k = a.len();
    for i in 0..b.len() {
        let ch = b[i];
        *window.entry(ch).or_insert(0) += 1;
        if i >= k {
            let left_ch = b[i - k];
            if let Some(v) = window.get_mut(&left_ch) {
                *v -= 1;
                if *v == 0 {
                    window.remove(&left_ch);
                }
            }
        }
        if window == need {
            return true;
        }
    }
    false
}
