fn search(nums: Vec<i64>, target: i64) -> i64 {
    let mut lo: i64 = 0;
    let mut hi: i64 = nums.len() as i64 - 1;
    while lo <= hi {
        let mid = (lo + hi) / 2;
        let v = nums[mid as usize];
        if v == target {
            return mid;
        } else if v < target {
            lo = mid + 1;
        } else {
            hi = mid - 1;
        }
    }
    -1
}
