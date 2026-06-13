use std::collections::HashSet;

fn can_partition(nums: Vec<i64>) -> bool {
    let total: i64 = nums.iter().sum();
    if total % 2 != 0 {
        return false;
    }
    let target = total / 2;
    let mut dp: HashSet<i64> = HashSet::new();
    dp.insert(0);
    for &n in nums.iter() {
        let mut next = dp.clone();
        for &s in dp.iter() {
            next.insert(s + n);
        }
        dp = next;
        if dp.contains(&target) {
            return true;
        }
    }
    false
}
