public class Solution {
    public bool can_jump(long[] nums) {
        long reach = 0;
        for (int i = 0; i < nums.Length; i++) {
            if (i > reach) return false;
            reach = Math.Max(reach, i + nums[i]);
        }
        return true;
    }
}
