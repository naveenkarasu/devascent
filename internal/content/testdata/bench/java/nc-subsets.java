import java.util.*;

class Solution {
    public long[][] subsets(long[] nums) {
        List<long[]> res = new ArrayList<>();
        backtrack(nums, 0, new ArrayList<>(), res);
        res.sort((a, b) -> {
            int len = Math.min(a.length, b.length);
            for (int i = 0; i < len; i++) {
                if (a[i] != b[i]) return Long.compare(a[i], b[i]);
            }
            return a.length - b.length;
        });
        return res.toArray(new long[0][]);
    }

    private void backtrack(long[] nums, int start, List<Long> current, List<long[]> res) {
        long[] sub = new long[current.size()];
        for (int i = 0; i < current.size(); i++) sub[i] = current.get(i);
        Arrays.sort(sub);
        res.add(sub);
        for (int i = start; i < nums.length; i++) {
            current.add(nums[i]);
            backtrack(nums, i + 1, current, res);
            current.remove(current.size() - 1);
        }
    }
}
