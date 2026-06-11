fn missing_number(nums: Vec<i64>) -> i64 {
    let n = nums.len() as i64;
    let expected = n * (n + 1) / 2;
    let sum: i64 = nums.iter().sum();
    expected - sum
}
