use std::collections::HashSet;

fn contains_duplicate(nums: Vec<i64>) -> bool {
    let mut seen: HashSet<i64> = HashSet::new();
    for n in nums {
        if seen.contains(&n) {
            return true;
        }
        seen.insert(n);
    }
    false
}
