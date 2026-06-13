public class Solution {
    public long[] next_permutation(long[] nums) {
        int n = nums.Length;
        int pivot = -1;
        for (int i = n - 2; i >= 0; i--) {
            if (nums[i] < nums[i + 1]) { pivot = i; break; }
        }
        if (pivot == -1) {
            System.Array.Reverse(nums);
            return nums;
        }
        for (int r = n - 1; r > pivot; r--) {
            if (nums[r] > nums[pivot]) {
                long tmp = nums[pivot]; nums[pivot] = nums[r]; nums[r] = tmp;
                break;
            }
        }
        int lo = pivot + 1, hi = n - 1;
        while (lo < hi) {
            long tmp = nums[lo]; nums[lo] = nums[hi]; nums[hi] = tmp;
            lo++; hi--;
        }
        return nums;
    }
}
