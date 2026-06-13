import java.util.*;

class Solution {
    public long[] next_permutation(long[] nums) {
        long[] arr = Arrays.copyOf(nums, nums.length);
        int n = arr.length;
        int i = n - 2;
        while (i >= 0 && arr[i] >= arr[i + 1]) i--;
        if (i != -1) {
            int j = n - 1;
            while (arr[j] <= arr[i]) j--;
            long tmp = arr[i]; arr[i] = arr[j]; arr[j] = tmp;
        }
        int left = i + 1, right = n - 1;
        while (left < right) {
            long tmp = arr[left]; arr[left] = arr[right]; arr[right] = tmp;
            left++;
            right--;
        }
        return arr;
    }
}
