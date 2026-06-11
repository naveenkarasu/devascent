use std::collections::HashMap;

fn least_interval(tasks: Vec<String>, n: i64) -> i64 {
    let mut counts: HashMap<String, i64> = HashMap::new();
    for t in &tasks {
        *counts.entry(t.clone()).or_insert(0) += 1;
    }
    let max_count = counts.values().cloned().max().unwrap_or(0);
    let num_max = counts.values().filter(|&&v| v == max_count).count() as i64;
    let candidate = (max_count - 1) * (n + 1) + num_max;
    let total = tasks.len() as i64;
    if total > candidate { total } else { candidate }
}
