use std::collections::HashMap;

fn find_unique_among_k(k: i64, arr: Vec<i64>) -> i64 {
    let mut counts: HashMap<i64, i64> = HashMap::new();
    let mut order: Vec<i64> = Vec::new();
    for &x in &arr {
        let e = counts.entry(x).or_insert(0);
        if *e == 0 {
            order.push(x);
        }
        *e += 1;
    }
    for val in order {
        if counts[&val] % k != 0 {
            return val;
        }
    }
    -1
}
