class Solution {
    public long first_missing_positive(long[] nums) {
        int n = nums.length;
        for (int i = 0; i < n; i++) {
            while (nums[i] >= 1 && nums[i] <= n && nums[(int) nums[i] - 1] != nums[i]) {
                int correct = (int) nums[i] - 1;
                long tmp = nums[i];
                nums[i] = nums[correct];
                nums[correct] = tmp;
            }
        }
        for (int i = 0; i < n; i++) {
            if (nums[i] != i + 1) return i + 1;
        }
        return n + 1;
    }
}
