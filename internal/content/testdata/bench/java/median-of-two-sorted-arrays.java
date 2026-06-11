class Solution {
    public double find_median_sorted_arrays(long[] nums1, long[] nums2) {
        if (nums1.length > nums2.length) {
            long[] tmp = nums1; nums1 = nums2; nums2 = tmp;
        }
        int m = nums1.length, n = nums2.length;
        int total = m + n;
        int lo = 0, hi = m;
        while (lo <= hi) {
            int cut1 = (lo + hi) / 2;
            int cut2 = (total + 1) / 2 - cut1;
            long l1 = (cut1 > 0) ? nums1[cut1 - 1] : Long.MIN_VALUE;
            long l2 = (cut2 > 0) ? nums2[cut2 - 1] : Long.MIN_VALUE;
            long r1 = (cut1 < m) ? nums1[cut1] : Long.MAX_VALUE;
            long r2 = (cut2 < n) ? nums2[cut2] : Long.MAX_VALUE;
            if (l1 <= r2 && l2 <= r1) {
                if (total % 2 == 0) {
                    return (Math.max(l1, l2) + Math.min(r1, r2)) / 2.0;
                } else {
                    return (double) Math.max(l1, l2);
                }
            } else if (l1 > r2) {
                hi = cut1 - 1;
            } else {
                lo = cut1 + 1;
            }
        }
        return 0.0;
    }
}
