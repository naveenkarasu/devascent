public class Solution {
    public long jump(long[] nums) {
        long jumps = 0;
        long curEnd = 0;
        long farthest = 0;
        for (int i = 0; i < nums.Length - 1; i++) {
            farthest = Math.Max(farthest, i + nums[i]);
            if (i == curEnd) {
                jumps++;
                curEnd = farthest;
            }
        }
        return jumps;
    }
}
