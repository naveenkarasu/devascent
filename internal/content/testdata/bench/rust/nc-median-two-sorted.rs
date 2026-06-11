fn find_median_sorted_arrays(nums1: Vec<i64>, nums2: Vec<i64>) -> i64 {
    let (a, b) = if nums1.len() > nums2.len() {
        (nums2, nums1)
    } else {
        (nums1, nums2)
    };
    let m = a.len();
    let n = b.len();
    let half = (m + n) / 2;
    let mut lo: i64 = 0;
    let mut hi: i64 = m as i64;
    loop {
        let i = (lo + hi) / 2; // partition index in a
        let j = half as i64 - i; // partition index in b
        let iu = i as usize;
        let ju = j as usize;
        let left1: i64 = if i > 0 { a[iu - 1] } else { i64::MIN };
        let right1: i64 = if (i as usize) < m { a[iu] } else { i64::MAX };
        let left2: i64 = if j > 0 { b[ju - 1] } else { i64::MIN };
        let right2: i64 = if (j as usize) < n { b[ju] } else { i64::MAX };
        if left1 <= right2 && left2 <= right1 {
            if (m + n) % 2 == 1 {
                return right1.min(right2);
            }
            return (left1.max(left2) + right1.min(right2)) / 2;
        } else if left1 > right2 {
            hi = i - 1;
        } else {
            lo = i + 1;
        }
    }
}
