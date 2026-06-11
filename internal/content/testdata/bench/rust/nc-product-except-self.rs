fn product_except_self(nums: Vec<i64>) -> Vec<i64> {
    let n = nums.len();
    let mut res = vec![1i64; n];
    let mut prefix = 1i64;
    for i in 0..n {
        res[i] = prefix;
        prefix *= nums[i];
    }
    let mut suffix = 1i64;
    for i in (0..n).rev() {
        res[i] *= suffix;
        suffix *= nums[i];
    }
    res
}
