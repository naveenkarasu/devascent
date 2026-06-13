fn majority_element(nums: Vec<i64>) -> i64 {
    let mut count = 0i64;
    let mut candidate = 0i64;
    for &n in &nums {
        if count == 0 {
            candidate = n;
        }
        count += if n == candidate { 1 } else { -1 };
    }
    let occ = nums.iter().filter(|&&x| x == candidate).count();
    if occ > nums.len() / 2 {
        candidate
    } else {
        -1
    }
}
