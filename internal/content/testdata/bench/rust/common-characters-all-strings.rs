use std::collections::HashSet;

fn count_common_characters(strings: Vec<String>) -> i64 {
    if strings.is_empty() {
        return 0;
    }
    let mut common: HashSet<char> = strings[0].chars().collect();
    for s in &strings[1..] {
        let cur: HashSet<char> = s.chars().collect();
        common = common.intersection(&cur).copied().collect();
    }
    common.len() as i64
}
