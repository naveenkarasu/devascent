import java.util.*;

class Solution {
    public long[] merge_sorted_in_place(long[] arr1, long[] arr2) {
        int m = arr1.length, n = arr2.length;
        long[] result = new long[m + n];
        int l = m - 1, r = n - 1, k = m + n - 1;
        while (l >= 0 && r >= 0) {
            if (arr1[l] > arr2[r]) {
                result[k--] = arr1[l--];
            } else {
                result[k--] = arr2[r--];
            }
        }
        while (l >= 0) result[k--] = arr1[l--];
        while (r >= 0) result[k--] = arr2[r--];
        return result;
    }
}
