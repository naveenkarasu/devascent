import java.util.*;

class Solution {
    public long trap_water(long[] height) {
        int l = 0, r = height.length - 1;
        long level = 0, water = 0;
        while (l < r) {
            long lower = (height[l] < height[r]) ? height[l] : height[r];
            if (height[l] < height[r]) l++;
            else r--;
            level = Math.max(level, lower);
            water += level - lower;
        }
        return water;
    }
}
