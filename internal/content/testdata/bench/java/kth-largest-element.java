import java.util.*;

class Solution {
    public long kth_largest(long[] nums, long k) {
        long[] sorted = Arrays.copyOf(nums, nums.length);
        Arrays.sort(sorted);
        // reverse: largest at index 0
        int left = 0, right = sorted.length - 1;
        while (left < right) {
            long tmp = sorted[left]; sorted[left] = sorted[right]; sorted[right] = tmp;
            left++; right--;
        }
        return sorted[(int)(k - 1)];
    }
}
