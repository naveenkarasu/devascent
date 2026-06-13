public class Solution {
    public double find_median_sorted_arrays(long[] nums1, long[] nums2) {
        if (nums1.Length > nums2.Length) {
            long[] tmp = nums1; nums1 = nums2; nums2 = tmp;
        }
        int m = nums1.Length, n = nums2.Length;
        int total = m + n;
        int lo = 0, hi = m;
        while (lo <= hi) {
            int cut1 = (lo + hi) / 2;
            int cut2 = (total + 1) / 2 - cut1;
            long l1 = (cut1 > 0) ? nums1[cut1 - 1] : long.MinValue;
            long l2 = (cut2 > 0) ? nums2[cut2 - 1] : long.MinValue;
            long r1 = (cut1 < m) ? nums1[cut1] : long.MaxValue;
            long r2 = (cut2 < n) ? nums2[cut2] : long.MaxValue;
            if (l1 <= r2 && l2 <= r1) {
                if (total % 2 == 0) {
                    return (Math.Max(l1, l2) + Math.Min(r1, r2)) / 2.0;
                } else {
                    return (double)Math.Max(l1, l2);
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
