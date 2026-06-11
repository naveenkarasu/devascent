fn max_partition_sum(arr: Vec<i64>) -> i64 {
    let total: i64 = arr.iter().sum();
    let mut current: i64 = 0;
    let mut ans: i64 = 0;
    for i in 0..arr.len().saturating_sub(1) {
        current += arr[i];
        ans = ans.max(current.max(total - current));
    }
    ans
}
