use std::collections::HashMap;

fn partition_labels(s: String) -> Vec<i64> {
    let chars: Vec<char> = s.chars().collect();
    let mut last: HashMap<char, i64> = HashMap::new();
    for (i, &c) in chars.iter().enumerate() {
        last.insert(c, i as i64);
    }
    let mut result: Vec<i64> = Vec::new();
    let mut start = 0i64;
    let mut end = 0i64;
    for (i, &c) in chars.iter().enumerate() {
        let l = *last.get(&c).unwrap();
        if l > end {
            end = l;
        }
        if i as i64 == end {
            result.push(end - start + 1);
            start = i as i64 + 1;
        }
    }
    result
}
