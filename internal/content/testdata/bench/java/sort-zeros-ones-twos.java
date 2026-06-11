import java.util.*;

class Solution {
    public long[] sort_012(long[] nums) {
        long[] arr = Arrays.copyOf(nums, nums.length);
        int low = 0, mid = 0, high = arr.length - 1;
        while (mid <= high) {
            if (arr[mid] == 0) {
                long tmp = arr[low]; arr[low] = arr[mid]; arr[mid] = tmp;
                low++; mid++;
            } else if (arr[mid] == 1) {
                mid++;
            } else {
                long tmp = arr[mid]; arr[mid] = arr[high]; arr[high] = tmp;
                high--;
            }
        }
        return arr;
    }
}
