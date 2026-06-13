fn find_majority_element(nums: Vec<i64>) -> i64 {
    let mut candidate: i64 = 0;
    let mut count: i64 = 0;
    for &x in &nums {
        if count == 0 {
            candidate = x;
        }
        if x == candidate {
            count += 1;
        } else {
            count -= 1;
        }
    }
    let freq = nums.iter().filter(|&&x| x == candidate).count() as i64;
    if freq > nums.len() as i64 / 2 {
        return candidate;
    }
    -1
}
