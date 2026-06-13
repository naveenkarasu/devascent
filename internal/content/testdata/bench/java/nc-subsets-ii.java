import java.util.*;

class Solution {
    public long[][] subsets_with_dup(long[] nums) {
        List<long[]> res = new ArrayList<>();
        Arrays.sort(nums);
        backtrack(nums, 0, new ArrayList<>(), res);
        // Sort result like Python: each subset sorted, then sort all subsets
        res.sort((a, b) -> {
            int len = Math.min(a.length, b.length);
            for (int i = 0; i < len; i++) {
                if (a[i] != b[i]) return Long.compare(a[i], b[i]);
            }
            return Integer.compare(a.length, b.length);
        });
        return res.toArray(new long[0][]);
    }

    private void backtrack(long[] nums, int start, List<Long> current, List<long[]> res) {
        long[] arr = new long[current.size()];
        for (int i = 0; i < current.size(); i++) arr[i] = current.get(i);
        Arrays.sort(arr);
        res.add(arr);
        for (int i = start; i < nums.length; i++) {
            if (i > start && nums[i] == nums[i - 1]) continue;
            current.add(nums[i]);
            backtrack(nums, i + 1, current, res);
            current.remove(current.size() - 1);
        }
    }
}
