fn find_median_sorted_arrays(nums1: Vec<i64>, nums2: Vec<i64>) -> f64 {
    let (a, b) = if nums1.len() > nums2.len() {
        (nums2, nums1)
    } else {
        (nums1, nums2)
    };
    let m = a.len();
    let n = b.len();
    let total = m + n;
    let mut lo: i64 = 0;
    let mut hi: i64 = m as i64;
    while lo <= hi {
        let cut1 = (lo + hi) / 2;
        let cut2 = (total as i64 + 1) / 2 - cut1;
        let l1 = if cut1 > 0 { a[(cut1 - 1) as usize] } else { i64::MIN };
        let l2 = if cut2 > 0 { b[(cut2 - 1) as usize] } else { i64::MIN };
        let r1 = if (cut1 as usize) < m { a[cut1 as usize] } else { i64::MAX };
        let r2 = if (cut2 as usize) < n { b[cut2 as usize] } else { i64::MAX };
        if l1 <= r2 && l2 <= r1 {
            if total % 2 == 0 {
                let max_l = l1.max(l2);
                let min_r = r1.min(r2);
                return (max_l + min_r) as f64 / 2.0;
            } else {
                return l1.max(l2) as f64;
            }
        } else if l1 > r2 {
            hi = cut1 - 1;
        } else {
            lo = cut1 + 1;
        }
    }
    0.0
}
