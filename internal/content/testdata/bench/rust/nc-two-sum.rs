use std::collections::HashMap;

fn two_sum(nums: Vec<i64>, target: i64) -> Vec<i64> {
    let mut seen: HashMap<i64, i64> = HashMap::new();
    for (i, &n) in nums.iter().enumerate() {
        let need = target - n;
        if let Some(&j) = seen.get(&need) {
            return vec![j, i as i64];
        }
        seen.insert(n, i as i64);
    }
    Vec::new()
}
