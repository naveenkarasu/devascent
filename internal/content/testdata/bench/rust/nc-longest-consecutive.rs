use std::collections::HashSet;

fn longest_consecutive(nums: Vec<i64>) -> i64 {
    let s: HashSet<i64> = nums.into_iter().collect();
    let mut best = 0i64;
    for &n in &s {
        if !s.contains(&(n - 1)) {
            let mut length = 1i64;
            while s.contains(&(n + length)) {
                length += 1;
            }
            best = best.max(length);
        }
    }
    best
}
