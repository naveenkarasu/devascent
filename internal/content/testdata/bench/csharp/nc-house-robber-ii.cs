public class Solution {
    private long RobLinear(long[] houses) {
        if (houses.Length == 0) return 0;
        if (houses.Length == 1) return houses[0];
        long prev2 = houses[0], prev1 = Math.Max(houses[0], houses[1]);
        for (int i = 2; i < houses.Length; i++) {
            long curr = Math.Max(prev1, prev2 + houses[i]);
            prev2 = prev1;
            prev1 = curr;
        }
        return prev1;
    }

    public long rob_circular(long[] nums) {
        if (nums.Length == 1) return nums[0];
        long[] a = new long[nums.Length - 1];
        long[] b = new long[nums.Length - 1];
        System.Array.Copy(nums, 0, a, 0, nums.Length - 1);
        System.Array.Copy(nums, 1, b, 0, nums.Length - 1);
        return Math.Max(RobLinear(a), RobLinear(b));
    }
}
