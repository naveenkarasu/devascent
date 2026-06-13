fn max_product(nums: Vec<i64>) -> i64 {
    let mut max_prod = nums[0];
    let mut min_prod = nums[0];
    let mut res = nums[0];
    for &n in nums.iter().skip(1) {
        let cands = [n, max_prod * n, min_prod * n];
        max_prod = *cands.iter().max().unwrap();
        min_prod = *cands.iter().min().unwrap();
        if max_prod > res {
            res = max_prod;
        }
    }
    res
}
