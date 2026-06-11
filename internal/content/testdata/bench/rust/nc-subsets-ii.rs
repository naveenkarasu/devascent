fn subsets_with_dup(nums: Vec<i64>) -> Vec<Vec<i64>> {
    let mut nums = nums;
    nums.sort();
    let mut res: Vec<Vec<i64>> = Vec::new();
    let mut current: Vec<i64> = Vec::new();

    fn backtrack(start: usize, nums: &Vec<i64>, current: &mut Vec<i64>, res: &mut Vec<Vec<i64>>) {
        res.push(current.clone());
        for i in start..nums.len() {
            if i > start && nums[i] == nums[i - 1] {
                continue;
            }
            current.push(nums[i]);
            backtrack(i + 1, nums, current, res);
            current.pop();
        }
    }

    backtrack(0, &nums, &mut current, &mut res);
    for s in res.iter_mut() {
        s.sort();
    }
    res.sort();
    res
}
