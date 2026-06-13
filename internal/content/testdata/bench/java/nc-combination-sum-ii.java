import java.util.*;

class Solution {
    public long[][] combination_sum2(long[] candidates, long target) {
        List<long[]> res = new ArrayList<>();
        Arrays.sort(candidates);
        backtrack(candidates, (int) target, 0, new ArrayList<>(), res);
        res.sort((a, b) -> {
            int len = Math.min(a.length, b.length);
            for (int i = 0; i < len; i++) {
                if (a[i] != b[i]) return Long.compare(a[i], b[i]);
            }
            return Integer.compare(a.length, b.length);
        });
        return res.toArray(new long[0][]);
    }

    private void backtrack(long[] candidates, int remaining, int start, List<Long> current, List<long[]> res) {
        if (remaining == 0) {
            long[] arr = new long[current.size()];
            for (int i = 0; i < current.size(); i++) arr[i] = current.get(i);
            Arrays.sort(arr);
            res.add(arr);
            return;
        }
        for (int i = start; i < candidates.length; i++) {
            if (candidates[i] > remaining) break;
            if (i > start && candidates[i] == candidates[i - 1]) continue;
            current.add(candidates[i]);
            backtrack(candidates, (int)(remaining - candidates[i]), i + 1, current, res);
            current.remove(current.size() - 1);
        }
    }
}
