use std::collections::HashMap;

fn find_target_sum_ways(nums: Vec<i64>, target: i64) -> i64 {
    let mut dp: HashMap<i64, i64> = HashMap::new();
    dp.insert(0, 1);
    for &num in nums.iter() {
        let mut next: HashMap<i64, i64> = HashMap::new();
        for (&s, &cnt) in dp.iter() {
            *next.entry(s + num).or_insert(0) += cnt;
            *next.entry(s - num).or_insert(0) += cnt;
        }
        dp = next;
    }
    *dp.get(&target).unwrap_or(&0)
}
