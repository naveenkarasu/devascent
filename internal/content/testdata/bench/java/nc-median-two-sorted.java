class Solution {
    public long find_median_sorted_arrays(long[] nums1, long[] nums2) {
        if (nums1.length > nums2.length) {
            long[] tmp = nums1; nums1 = nums2; nums2 = tmp;
        }
        int m = nums1.length, n = nums2.length;
        int half = (m + n) / 2;
        int lo = 0, hi = m;
        while (true) {
            int i = (lo + hi) / 2;
            int j = half - i;
            long left1  = i > 0 ? nums1[i - 1] : Long.MIN_VALUE;
            long right1 = i < m ? nums1[i]     : Long.MAX_VALUE;
            long left2  = j > 0 ? nums2[j - 1] : Long.MIN_VALUE;
            long right2 = j < n ? nums2[j]     : Long.MAX_VALUE;
            if (left1 <= right2 && left2 <= right1) {
                if ((m + n) % 2 == 1) return Math.min(right1, right2);
                return (Math.max(left1, left2) + Math.min(right1, right2)) / 2;
            } else if (left1 > right2) {
                hi = i - 1;
            } else {
                lo = i + 1;
            }
        }
    }
}
