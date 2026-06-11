import java.util.*;

class Solution {
    public long jump(long[] nums) {
        long jumps = 0;
        long curEnd = 0;
        long farthest = 0;
        for (int i = 0; i < nums.length - 1; i++) {
            farthest = Math.max(farthest, i + nums[i]);
            if (i == curEnd) {
                jumps++;
                curEnd = farthest;
            }
        }
        return jumps;
    }
}
