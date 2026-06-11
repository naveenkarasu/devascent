import java.util.*;

class Solution {
    public long[] merge_sorted_arrays(long[] arr1, long[] arr2) {
        long[] result = new long[arr1.length + arr2.length];
        int i = 0, j = 0, k = 0;
        while (i < arr1.length && j < arr2.length) {
            if (arr1[i] <= arr2[j]) result[k++] = arr1[i++];
            else result[k++] = arr2[j++];
        }
        while (i < arr1.length) result[k++] = arr1[i++];
        while (j < arr2.length) result[k++] = arr2[j++];
        return result;
    }
}
