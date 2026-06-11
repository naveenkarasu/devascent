function find_median_sorted_arrays(nums1: number[], nums2: number[]): number {
    if (nums1.length > nums2.length) {
        [nums1, nums2] = [nums2, nums1];
    }
    const m = nums1.length, n = nums2.length;
    const total = m + n;
    let lo = 0, hi = m;
    while (lo <= hi) {
        const cut1 = Math.floor((lo + hi) / 2);
        const cut2 = Math.floor((total + 1) / 2) - cut1;
        const l1 = cut1 > 0 ? nums1[cut1 - 1] : -Infinity;
        const l2 = cut2 > 0 ? nums2[cut2 - 1] : -Infinity;
        const r1 = cut1 < m ? nums1[cut1] : Infinity;
        const r2 = cut2 < n ? nums2[cut2] : Infinity;
        if (l1 <= r2 && l2 <= r1) {
            if (total % 2 === 0) {
                return (Math.max(l1, l2) + Math.min(r1, r2)) / 2.0;
            } else {
                return Math.max(l1, l2);
            }
        } else if (l1 > r2) {
            hi = cut1 - 1;
        } else {
            lo = cut1 + 1;
        }
    }
    return 0.0;
}
