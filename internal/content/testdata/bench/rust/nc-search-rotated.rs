fn search_rotated(nums: Vec<i64>, target: i64) -> i64 {
    let mut lo: i64 = 0;
    let mut hi: i64 = nums.len() as i64 - 1;
    while lo <= hi {
        let mid = (lo + hi) / 2;
        let vmid = nums[mid as usize];
        if vmid == target {
            return mid;
        }
        let vlo = nums[lo as usize];
        let vhi = nums[hi as usize];
        if vlo <= vmid {
            // left half is sorted
            if vlo <= target && target < vmid {
                hi = mid - 1;
            } else {
                lo = mid + 1;
            }
        } else {
            // right half is sorted
            if vmid < target && target <= vhi {
                lo = mid + 1;
            } else {
                hi = mid - 1;
            }
        }
    }
    -1
}
