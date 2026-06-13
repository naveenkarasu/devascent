import java.util.*;

class Solution {
    public long[][] three_sum(long[] nums) {
        Arrays.sort(nums);
        List<long[]> res = new ArrayList<>();
        int n = nums.length;
        for (int i = 0; i < n; i++) {
            if (i > 0 && nums[i] == nums[i - 1]) continue;
            int lo = i + 1, hi = n - 1;
            while (lo < hi) {
                long total = nums[i] + nums[lo] + nums[hi];
                if (total < 0) lo++;
                else if (total > 0) hi--;
                else {
                    res.add(new long[]{nums[i], nums[lo], nums[hi]});
                    lo++;
                    hi--;
                    while (lo < hi && nums[lo] == nums[lo - 1]) lo++;
                    while (lo < hi && nums[hi] == nums[hi + 1]) hi--;
                }
            }
        }
        return res.toArray(new long[0][]);
    }
}
