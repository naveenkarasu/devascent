import java.util.*;

class Solution {
    public long[][] permute(long[] nums) {
        List<long[]> res = new ArrayList<>();
        backtrack(nums, new boolean[nums.length], new ArrayList<>(), res);
        res.sort((a, b) -> {
            for (int i = 0; i < a.length; i++) {
                if (a[i] != b[i]) return Long.compare(a[i], b[i]);
            }
            return 0;
        });
        return res.toArray(new long[0][]);
    }

    private void backtrack(long[] nums, boolean[] used, List<Long> current, List<long[]> res) {
        if (current.size() == nums.length) {
            long[] perm = new long[nums.length];
            for (int i = 0; i < nums.length; i++) perm[i] = current.get(i);
            res.add(perm);
            return;
        }
        for (int i = 0; i < nums.length; i++) {
            if (!used[i]) {
                used[i] = true;
                current.add(nums[i]);
                backtrack(nums, used, current, res);
                current.remove(current.size() - 1);
                used[i] = false;
            }
        }
    }
}
