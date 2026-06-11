public class Solution {
    public long find_median_sorted_arrays(long[] nums1, long[] nums2) {
        if (nums1.Length > nums2.Length) {
            var tmp = nums1; nums1 = nums2; nums2 = tmp;
        }
        int m = nums1.Length, n = nums2.Length;
        int half = (m + n) / 2;
        int lo = 0, hi = m;
        while (true) {
            int i = (lo + hi) / 2;
            int j = half - i;
            long left1  = i > 0 ? nums1[i - 1] : long.MinValue;
            long right1 = i < m ? nums1[i]     : long.MaxValue;
            long left2  = j > 0 ? nums2[j - 1] : long.MinValue;
            long right2 = j < n ? nums2[j]     : long.MaxValue;
            if (left1 <= right2 && left2 <= right1) {
                if ((m + n) % 2 == 1)
                    return Math.Min(right1, right2);
                return (long)((Math.Max(left1, left2) + Math.Min(right1, right2)) / 2.0);
            } else if (left1 > right2) {
                hi = i - 1;
            } else {
                lo = i + 1;
            }
        }
    }
}
