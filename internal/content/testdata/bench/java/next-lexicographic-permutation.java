import java.util.*;

class Solution {
    public long[] next_permutation(long[] nums) {
        int n = nums.length;
        int pivot = -1;
        for (int i = n - 2; i >= 0; i--) {
            if (nums[i] < nums[i + 1]) {
                pivot = i;
                break;
            }
        }
        if (pivot == -1) {
            // reverse the whole array
            for (int i = 0, j = n - 1; i < j; i++, j--) {
                long tmp = nums[i]; nums[i] = nums[j]; nums[j] = tmp;
            }
            return nums;
        }
        // find rightmost element greater than pivot
        for (int r = n - 1; r > pivot; r--) {
            if (nums[r] > nums[pivot]) {
                long tmp = nums[r]; nums[r] = nums[pivot]; nums[pivot] = tmp;
                break;
            }
        }
        // reverse from pivot+1 to end
        for (int i = pivot + 1, j = n - 1; i < j; i++, j--) {
            long tmp = nums[i]; nums[i] = nums[j]; nums[j] = tmp;
        }
        return nums;
    }
}
