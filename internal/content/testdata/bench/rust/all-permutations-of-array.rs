fn all_permutations(nums: Vec<i64>) -> Vec<Vec<i64>> {
    if nums.is_empty() {
        return vec![vec![]];
    }
    let mut result: Vec<Vec<i64>> = Vec::new();
    for i in 0..nums.len() {
        let n = nums[i];
        let mut rest: Vec<i64> = Vec::new();
        rest.extend_from_slice(&nums[..i]);
        rest.extend_from_slice(&nums[i + 1..]);
        for perm in all_permutations(rest) {
            let mut combined = vec![n];
            combined.extend(perm);
            result.push(combined);
        }
    }
    result
}
