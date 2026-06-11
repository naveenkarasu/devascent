class Solution {
    private long robLinear(long[] houses) {
        if (houses.length == 0) return 0;
        if (houses.length == 1) return houses[0];
        long prev2 = houses[0], prev1 = Math.max(houses[0], houses[1]);
        for (int i = 2; i < houses.length; i++) {
            long curr = Math.max(prev1, prev2 + houses[i]);
            prev2 = prev1;
            prev1 = curr;
        }
        return prev1;
    }

    public long rob_circular(long[] nums) {
        if (nums.length == 1) return nums[0];
        long[] sub1 = java.util.Arrays.copyOfRange(nums, 0, nums.length - 1);
        long[] sub2 = java.util.Arrays.copyOfRange(nums, 1, nums.length);
        return Math.max(robLinear(sub1), robLinear(sub2));
    }
}
